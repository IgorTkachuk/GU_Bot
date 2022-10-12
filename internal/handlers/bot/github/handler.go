package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/NEKETSKY/footg-bot/internal/bot"
	config2 "github.com/NEKETSKY/footg-bot/internal/config"
	"github.com/NEKETSKY/footg-bot/internal/domain/bothandler"
	"github.com/NEKETSKY/footg-bot/internal/domain/github"
	"github.com/NEKETSKY/footg-bot/internal/websvcclient"
	tele "gopkg.in/telebot.v3"
	"io/ioutil"
	"net/http"
	"time"
)

var _ bothandler.BotHandler = &Handler{}

const (
	getRepoInfoCmd = "/repoinfo"
	ghApiUrl       = "https://api.github.com"
)

type Handler struct {
	bot *bot.FooBot
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(b *bot.FooBot) {
	h.bot = b

	b.Bot.Handle(getRepoInfoCmd, h.GetRepoInfo)
}

func (h *Handler) GetRepoInfo(context tele.Context) error {
	// TODO Get GitHub repo info
	// 1. Get info about user active repo from Web service over REST or gRPC
	// 2. Make request to GitHub API for get README.md content

	// Suppose in this point we have repo name (omit step 1 in previous list)

	config := config2.GetConfig()

	c := http.Client{Timeout: time.Duration(1) * time.Second}

	wCli := &websvcclient.Client{}
	userRepo := wCli.GetUserRepo(nil)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/repos/%s/contents/README.md", ghApiUrl, userRepo), nil)
	if err != nil {
		fmt.Printf("error %s", err)
		return err
	}
	req.Header.Add("Accept", `application/vnd.github+json`)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.GithubToken))

	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error %s", err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error %s", err)
		return err
	}

	var ghResp github.Response
	err = json.Unmarshal(body, &ghResp)
	if err != nil {
		return err
	}
	decodedBody, err := base64.StdEncoding.DecodeString(ghResp.Content)
	if err != nil {
		return err
	}

	err = context.Send(string(decodedBody))
	if err != nil {
		return err
	}

	return nil
}
