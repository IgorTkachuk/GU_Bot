package events

import (
	"encoding/json"
	"github.com/NEKETSKY/footg-bot/internal/bot"
	"github.com/NEKETSKY/footg-bot/internal/domain/restdto"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

const (
	newsUrl = "/api/events/:groupname"
)

type handler struct {
	bot *bot.FooBot
}

func NewHandler(bot *bot.FooBot) *handler {
	return &handler{bot: bot}
}

func (h *handler) Register(r *httprouter.Router) {
	r.HandlerFunc(http.MethodPost, newsUrl, h.Events)
}

func (h *handler) Events(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	groupName := params.ByName("groupname")

	var dto restdto.EventDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		log.Println("Can't decode json from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		return
	}
	h.bot.SendUnicast(groupName, dto.Message)
	w.WriteHeader(http.StatusOK)
}
