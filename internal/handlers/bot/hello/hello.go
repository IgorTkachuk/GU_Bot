package hello

import (
	"fmt"
	"github.com/NEKETSKY/footg-bot/internal/bot"
	tele "gopkg.in/telebot.v3"
	"log"
)

const (
	helloCmd = "/hello"
)

type Handler struct {
	bot *bot.FooBot
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(b *bot.FooBot) {
	h.bot = b
	b.Bot.Handle(helloCmd, h.ProcessHello)
}

func (h Handler) ProcessHello(c tele.Context) error {
	var (
		user = c.Sender()
	)
	log.Println(user.Recipient())

	return c.Send(fmt.Sprintf("Hello, %s", user.Username))
}
