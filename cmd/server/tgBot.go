package main

import (
	"context"

	logger "gitlab.ozon.dev/pircuser61/catalog/internal/logger"
	goodPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/core/good"
	botPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/tgbot"
	cmdAddPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/tgbot/command/add"
	cmdDeletePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/tgbot/command/delete"
	cmdGetPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/tgbot/command/get"
	cmdHelpPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/tgbot/command/help"
	cmdListPkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/tgbot/command/list"
	cmdUpdatePkg "gitlab.ozon.dev/pircuser61/catalog/internal/pkg/transport/tgbot/command/update"
	"go.uber.org/zap"
)

func runBot(ctx context.Context, good goodPkg.Repository) {

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

		commandHelp := cmdHelpPkg.New(bot.GetCmdList())
		bot.RgiesterHandler(commandHelp)
	}

	if err := bot.Run(ctx); err != nil {
		logger.Panic("bot.Run", zap.Error(err))
	}
}
