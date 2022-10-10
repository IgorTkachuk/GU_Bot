package main

import (
	"fmt"
	"github.com/NEKETSKY/footg-bot/internal/bot"
	"github.com/NEKETSKY/footg-bot/internal/config"
	"github.com/NEKETSKY/footg-bot/internal/handlers/bot/global"
	"github.com/NEKETSKY/footg-bot/internal/handlers/bot/hello"
	registration2 "github.com/NEKETSKY/footg-bot/internal/handlers/bot/registration"
	"github.com/NEKETSKY/footg-bot/internal/handlers/rest/events"
	"github.com/NEKETSKY/footg-bot/internal/handlers/rest/news"
	bot2 "github.com/NEKETSKY/footg-bot/internal/middleware/bot"
	"github.com/go-redis/redis/v9"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	config := config.GetConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.Db,
	})

	b, err := bot.NewFooBot(config.Token, rdb)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.RegisterMiddleware(bot2.SyncUserMiddleware)

	helloHandler := hello.NewHandler()
	helloHandler.Register(b)

	registrationHandler := registration2.NewHandler()
	registrationHandler.Register(b)

	globalHandler := global.NewHandler(registrationHandler)
	globalHandler.Register(b)

	log.Println("Trying FooTG-Bot starting for service")
	go b.Bot.Start()

	router := httprouter.New()

	restNewsHandler := news.NewHandler(b)
	restNewsHandler.Register(router)

	restEvensHandler := events.NewHandler(b)
	restEvensHandler.Register(router)

	log.Println("Trying REST API starting for service")
	log.Fatal(http.ListenAndServe(":8080", router))
}
