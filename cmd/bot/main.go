package main

import (
	"log"

	"gitlab.ozon.dev/pircuser61/catalog/internal/commander"
	"gitlab.ozon.dev/pircuser61/catalog/internal/handlers"
)

func main() {
	cmd, err := commander.Init()
	if err != nil {
		log.Panic(err)
	}
	handlers.Register(cmd)
	cmd.Run()
}
