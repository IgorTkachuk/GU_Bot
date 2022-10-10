package news

import (
	"github.com/NEKETSKY/footg-bot/internal/bot"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

const (
	newsUrl = "/api/news"
)

type handler struct {
	bot *bot.FooBot
}

func NewHandler(bot *bot.FooBot) *handler {
	return &handler{bot: bot}
}

func (h *handler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodPost, newsUrl, h.News)
}

func (h *handler) News(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	message, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	h.bot.SendBroadcast(string(message))
}
