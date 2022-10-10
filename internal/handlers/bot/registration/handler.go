package registration

import (
	"fmt"
	"github.com/NEKETSKY/footg-bot/internal/bot"
	"github.com/NEKETSKY/footg-bot/internal/domain/bothandler"
	"github.com/NEKETSKY/footg-bot/internal/fsm/registration"
	"github.com/NEKETSKY/footg-bot/pkg/cache"
	"github.com/NEKETSKY/footg-bot/pkg/cache/freecache"
	"github.com/enescakir/emoji"
	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
)

var _ bothandler.BotHandler = &Handler{}

var (
	menu = &tele.ReplyMarkup{
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	// Reply buttons.
	btnReplay = menu.Text(fmt.Sprintf("%v Replay", emoji.CounterclockwiseArrowsButton))
	btnCancel = menu.Text(fmt.Sprintf("%v Cancel", emoji.Prohibited))

	resetMenu = &tele.ReplyMarkup{
		RemoveKeyboard: true,
	}
)

const (
	registerCmd = "/register"
)

type Handler struct {
	bot        *bot.FooBot
	cacheFSM   cache.Repository
	cacheToken cache.Repository
}

func NewHandler() *Handler {
	cacheFSM := freecache.NewCacheRepo(10)
	cacheToken := freecache.NewCacheRepo(10)
	return &Handler{
		cacheFSM:   cacheFSM,
		cacheToken: cacheToken,
	}
}

func (h *Handler) Register(b *bot.FooBot) {
	h.bot = b

	b.Bot.Handle(&btnCancel, h.Cancel)
	b.Bot.Handle(&btnReplay, h.Replay)

	b.Bot.Handle(registerCmd, h.StartRegister)
}

func (h *Handler) StartRegister(c tele.Context) error {
	var (
		user = c.Sender()
	)

	// TODO is user in store?!!
	u := h.bot.GetUserInfo(user.Username)
	if u.Registered {
		c.Send("You already registered")
		return nil
	}

	h.bot.SessionState.Set([]byte(user.Username), []byte(bot.SessionRegister), 0)

	fsm := registration.NewRegisterFSM()
	err := fsm.FSM.Event(registration.EventRegister)
	if err != nil {
		return err
	}

	h.cacheFSM.Set([]byte(user.Username), []byte(fsm.FSM.Current()), 0)

	userToken := uuid.New()
	fmt.Printf("User %s token is: %s\n", user.Username, userToken)
	h.cacheToken.Set([]byte(user.Username), []byte(userToken.String()), 0)

	return c.Send("Enter you token from WeareComunity profile", resetMenu)
}

func (h *Handler) ProcessRegister(c tele.Context) error {
	var (
		user = c.Sender()
		text = c.Text()
	)

	menu.Reply(
		menu.Row(btnCancel),
		menu.Row(btnReplay),
	)

	userTokenBytes, err := h.cacheToken.Get([]byte(user.Username))
	if err != nil {
		return err
	}

	userToken := string(userTokenBytes)
	var event string

	if text == userToken {
		event = registration.EventGotRightToken
	} else {
		event = registration.EventGotWrongToken
	}

	rState, err := h.cacheFSM.Get([]byte(user.Username))
	if err != nil {
		return err
	}

	fsm := registration.NewRegisterFSM()
	fsm.FSM.SetState(string(rState))

	err = fsm.FSM.Event(event)
	if err != nil {
		return err
	}

	if fsm.FSM.Current() == registration.StateFailure {
		c.Send("You provide wrong token", menu)
		if err != nil {
			return err
		}
	} else if fsm.FSM.Current() == registration.StateSuccess {
		c.Send(fmt.Sprintf("Registration successful! %v", emoji.PartyPopper), resetMenu)
		if err != nil {
			return err
		}
		h.bot.SessionState.Set([]byte(user.Username), []byte(bot.SessionIdle), 0)
		h.cacheFSM.Del([]byte(user.Username))

		h.bot.MakeUserRegistered(user.Username)
	}

	return nil
}

func (h *Handler) Cancel(c tele.Context) error {
	var (
		user = c.Sender()
	)
	h.bot.SessionState.Set([]byte(user.Username), []byte(bot.SessionIdle), 0)
	h.cacheFSM.Del([]byte(user.Username))
	return c.Send("Registration process canceled!", resetMenu)
}

func (h *Handler) Replay(c tele.Context) error {
	return c.Send("Enter you token from WeareComunity profile", resetMenu)
}
