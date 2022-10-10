package global

import (
	"github.com/NEKETSKY/footg-bot/internal/bot"
	"github.com/NEKETSKY/footg-bot/internal/domain/bothandler"
	"github.com/NEKETSKY/footg-bot/internal/handlers/bot/registration"
	tele "gopkg.in/telebot.v3"
)

var _ bothandler.BotHandler = &handler{}

type handler struct {
	bot      *bot.FooBot
	register *registration.Handler
}

func NewHandler(register *registration.Handler) *handler {
	return &handler{
		register: register,
	}
}

func (h handler) Register(b *bot.FooBot) {
	h.bot = b
	b.Bot.Handle(tele.OnText, h.OnText)
}

func (h handler) OnText(c tele.Context) error {
	var (
		user = c.Sender()
	)

	currentSessionState, err := h.bot.SessionState.Get([]byte(user.Username))
	if err != nil {
		return err
	}

	if bot.UserSessionState(currentSessionState) == bot.SessionRegister {
		err = h.register.ProcessRegister(c)
		if err != nil {
			return err
		}
	}

	return nil
}
