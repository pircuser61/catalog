package main

import (
	"log"

	botPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command/add"
	cmdHelpPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command/help"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
)

func main() {
	good := goodPkg.New()

	var bot botPkg.Interface
	{
		bot = botPkg.MustNew()
		commandAdd := cmdAddPkg.New(good)
		bot.RgiesterHandler((commandAdd))

		commandHelp := cmdHelpPkg.New(bot.GetCmdList())
		bot.RgiesterHandler(commandHelp)
	}

	if err := bot.Run(); err != nil {
		log.Panic(err)
	}
}
