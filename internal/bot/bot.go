package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NEKETSKY/footg-bot/internal/domain/botdto"
	"github.com/NEKETSKY/footg-bot/pkg/utils"
	"github.com/coocood/freecache"
	"github.com/go-redis/redis/v9"
	"gopkg.in/telebot.v3"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

type UserSessionState string

const (
	SessionIdle     UserSessionState = "idle"
	SessionRegister UserSessionState = "register"

	UserBroadcastGroupName = "all"
)

type FooBot struct {
	Bot          *telebot.Bot
	SessionState *freecache.Cache
	Redis        *redis.Client
}

func NewFooBot(token string, r *redis.Client) (*FooBot, error) {
	sessionStateCache := freecache.NewCache(10)
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &FooBot{Bot: b, SessionState: sessionStateCache, Redis: r}, nil
}

func (b *FooBot) SendBroadcast(message string) error {
	chats := b.getChatIdByGroupName(UserBroadcastGroupName)
	for _, chat := range chats {
		_, err := b.Bot.Send(chat, message)
		if err != nil {
			return err
			log.Println(err)
		}
	}

	return nil
}

func (b *FooBot) SendUnicast(groupName string, message string) error {
	chats := b.getChatIdByGroupName(groupName)
	for _, chat := range chats {
		_, err := b.Bot.Send(chat, message)
		if err != nil {
			return err
			log.Println(err)
		}
	}

	return nil
}

func (b *FooBot) getChatIdByGroupName(g string) []*tele.Chat {
	ctx := context.Background()
	keys, err := b.Redis.Keys(ctx, "*").Result()

	if err != nil {
		log.Println("Can`t read all users keys from store")
	}

	chats := make([]*tele.Chat, 0)
	for _, username := range keys {
		u, err := b.Redis.Get(ctx, username).Result()
		if err != nil {
			log.Println(fmt.Sprintf("Can`t read from store info for user %s", username))
		}

		var userInfo botdto.UserInfoDTO
		err = json.Unmarshal([]byte(u), &userInfo)
		if err != nil {
			log.Println("Failed to decode user info due retrieve groups information")
		}

		if !utils.Contains(userInfo.Groups, g) {
			continue
		}

		chat, err := b.Bot.ChatByUsername(userInfo.ChatId)
		if err != nil {
			log.Println(err)
		}
		chats = append(chats, chat)
	}

	return chats
}

func (b *FooBot) MakeUserRegistered(username string) {
	usrDTO := b.GetUserInfo(username)
	usrDTO.Registered = true

	b.SetUserInfo(username, *usrDTO)
}

func (b *FooBot) GetUserInfo(username string) *botdto.UserInfoDTO {
	ctx := context.Background()
	u, err := b.Redis.Get(ctx, username).Result()
	if err != nil {
		log.Println(fmt.Sprintf("Can`t read from store info for user %s", username))
		return nil
	}

	var usrDTO botdto.UserInfoDTO
	err = json.Unmarshal([]byte(u), &usrDTO)
	if err != nil {
		log.Println("Failed unmarshal user info due making user registered")
		return nil
	}

	return &usrDTO
}

func (b *FooBot) SetUserInfo(username string, usrDTO botdto.UserInfoDTO) error {
	ctx := context.Background()
	userbytes, err := json.Marshal(usrDTO)
	if err != nil {
		log.Println("Failed marshal user info due making user registered")
		return err
	}

	b.Redis.Set(ctx, username, string(userbytes), 0)
	return nil
}

type md func(c tele.Context, b *FooBot)

func (b *FooBot) RegisterMiddleware(middleware md) {
	b.Bot.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			middleware(c, b)
			return next(c)
		}
	})
}
