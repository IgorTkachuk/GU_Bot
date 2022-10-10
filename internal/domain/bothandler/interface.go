package bothandler

import "github.com/NEKETSKY/footg-bot/internal/bot"

type BotHandler interface {
	Register(b *bot.FooBot)
}
