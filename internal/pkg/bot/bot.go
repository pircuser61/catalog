package commander

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"

	configPkg "gitlab.ozon.dev/pircuser61/catalog/internal/config"
	commandPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command"
)

var UnknownCommand = errors.New("Unknown command, /help for help")

type commander struct {
	bot   *tgbotapi.BotAPI
	route map[string]commandPkg.Interface
}

type Interface interface {
	Run() error
	RgiesterHandler(f commandPkg.Interface)
	GetCmdList() map[string]string
}

func MustNew() Interface {
	cfg, err := configPkg.GetTgBotConfig()
	if err != nil {
		log.Panic(errors.Wrap(err, "Init tgbot config"))
	}

	bot, err := tgbotapi.NewBotAPI(cfg.ApiKey)
	if err != nil {
		log.Panic(errors.Wrap(err, "Init tgbot"))
	}
	bot.Debug = cfg.Debug
	return &commander{bot, make(map[string]commandPkg.Interface)}

}

// RegisterHandler not thread-safe
func (c *commander) RgiesterHandler(handler commandPkg.Interface) {
	c.route[handler.Name()] = handler
}

func (c *commander) Run() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updatesChannel := c.bot.GetUpdatesChan(updateConfig)
	for update := range updatesChannel {
		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if cmdName := update.Message.Command(); cmdName != "" {
			if handler, ok := c.route[cmdName]; ok {
				msg.Text = handler.Process(update.Message.CommandArguments())
			} else {
				msg.Text = UnknownCommand.Error()
			}
		} else {
			msg.Text = fmt.Sprintf("you send <%v>", update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
		}
		_, err := c.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "Send messaege error")
		}
	}
	return nil
}

func (c *commander) GetCmdList() map[string]string {
	result := map[string]string{}
	for name, cmd := range c.route {
		result[name] = cmd.Description()
	}
	return result
}
