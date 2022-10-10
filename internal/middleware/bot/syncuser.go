package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NEKETSKY/footg-bot/internal/bot"
	"github.com/NEKETSKY/footg-bot/internal/domain/botdto"
	tele "gopkg.in/telebot.v3"
	"log"
	"strconv"
)

func SyncUserMiddleware(c tele.Context, b *bot.FooBot) {
	var (
		user = c.Sender()
	)

	ctx := context.Background()
	userInfo, err := b.Redis.Get(ctx, user.Username).Result()
	if err != nil {
		log.Println(fmt.Sprintf("Username %s isn`t exist in store", user.Username))

		gs := make([]string, 0)
		gs = append(gs, bot.UserBroadcastGroupName)

		u := botdto.UserInfoDTO{
			Registered: false,
			Groups:     gs,
			ChatId:     strconv.Itoa(int(c.Chat().ID)),
		}
		uBytes, err := json.Marshal(u)
		if err != nil {
			log.Println("Can't marshal user info due user sync")
		}

		b.Redis.Set(ctx, user.Username, string(uBytes), 0)

		return
	}

	var u botdto.UserInfoDTO
	err = json.Unmarshal([]byte(userInfo), &u)
	if err != nil {
		log.Println("Can't unmarshal user info due user sync")
	}

	log.Println(fmt.Sprintf("User %s is alerady present in sore. His chat ID equal %s", user.Username, u.ChatId))
}
