package main

import (
	"context"
	"log"

	botPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command/add"
	cmdDeletePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command/delete"
	cmdGetPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command/get"
	cmdListPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command/list"
	cmdUpdatePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command/update"

	cmdHelpPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command/help"
	cmdLockPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command/lock"
	cmdStopPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/bot/command/stop"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
)

func runBot(ctx context.Context, stopCh chan<- struct{}, good goodPkg.Interface) {

	var bot botPkg.Interface
	{
		bot = botPkg.MustNew()
		commandAdd := cmdAddPkg.New(good)
		bot.RgiesterHandler((commandAdd))

		commandCreate := cmdUpdatePkg.New(good)
		bot.RgiesterHandler((commandCreate))

		commandGet := cmdGetPkg.New(good)
		bot.RgiesterHandler((commandGet))

		commandDelete := cmdDeletePkg.New(good)
		bot.RgiesterHandler((commandDelete))

		commandList := cmdListPkg.New(good)
		bot.RgiesterHandler((commandList))

		commandStop := cmdStopPkg.New(stopCh)
		bot.RgiesterHandler((commandStop))

		commandLock := cmdLockPkg.New(good)
		bot.RgiesterHandler(commandLock)

		commandHelp := cmdHelpPkg.New(bot.GetCmdList())
		bot.RgiesterHandler(commandHelp)
	}

	if err := bot.Run(ctx); err != nil {
		log.Panic(err)
	}
}
