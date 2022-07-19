package commander

import (
	"fmt"
	"go1/config"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

const (
	help = "help"
)

type Commander struct {
	bot *tgbotapi.BotAPI
}

type CmdHandler struct {
	f    func(string) string
	help string
}

var route map[string]CmdHandler
var UnknownCommand = errors.New("Unknown command, /help for help")

func Init() (*Commander, error) {
	route = make(map[string]CmdHandler)

	bot, err := tgbotapi.NewBotAPI(config.ApiKey)
	if err != nil {
		return nil, errors.Wrap(err, "Init tgbot")
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	cmr := &Commander{bot}
	cmr.Register(help, getHelp, "help")
	return cmr, nil
}

func (*Commander) Register(cmd string, f func(string) string, helpStr string) {
	route[cmd] = CmdHandler{f, helpStr}
}

func (cmd *Commander) Run() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updatesChannel := cmd.bot.GetUpdatesChan(updateConfig)
	for update := range updatesChannel {
		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message == nil {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if cmd := update.Message.Command(); cmd != "" {
			if x, ok := route[cmd]; ok {
				msg.Text = x.f(update.Message.CommandArguments())
			} else {
				msg.Text = UnknownCommand.Error()
			}
		} else {
			msg.Text = fmt.Sprintf("you send <%v>", update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
		}
		_, err := cmd.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "Send messaege error")
		}
	}
	return nil
}

func getHelp(_ string) string {
	res := make([]string, 0, len(route))
	for key, x := range route {

		if key != help {
			res = append(res, x.help)
		}

	}
	return strings.Join(res, "\n")
}
