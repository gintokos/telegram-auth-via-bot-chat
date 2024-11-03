package telegramauth

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// link is the link with which the hash is merged (https://somesite/)
type Bot struct {
	*tgbotapi.BotAPI
	ustore        *userTokenStore
	usertokensize int
	link          string
}

type BotOptions struct {
	TokenBot, Link string
	UserTokenSize  int
	TTLusertoken   time.Duration
}

func NewAuthBot(opt BotOptions) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(opt.TokenBot)
	if err != nil {
		return &Bot{}, err
	}

	ustore := newTokenStore(opt.TTLusertoken)

	return &Bot{
		BotAPI:        bot,
		link:          opt.Link,
		usertokensize: opt.UserTokenSize,
		ustore:        ustore,
	}, nil
}

func (b *Bot) Start() {
	log.Printf("Authorized on account %s", b.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "login":
				timenow := strconv.FormatInt(time.Now().Unix(), 10)
				str := strconv.FormatInt(update.Message.Chat.ID, 10) + timenow
				fmt.Println(str)

				utoken := getUserToken(str, b.usertokensize)
				b.ustore.SaveUserToken(utoken)

				loginlink := b.link + utoken

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, loginlink)
				b.Send(msg)
			case "validate":
				var msg tgbotapi.MessageConfig
				if b.IsUsertokenExists(update.Message.CommandArguments()) {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "true"+update.Message.CommandArguments())
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "false"+update.Message.CommandArguments())
				}
				b.Send(msg)
			}
		}
	}
}

func (b *Bot) IsUsertokenExists(usertoken string) bool {
	return b.ustore.ValidateUserToken(usertoken)
}
