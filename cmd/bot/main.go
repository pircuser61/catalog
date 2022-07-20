package main

import (
	"log"

	"Catalog/internal/commander"
	"Catalog/internal/handlers"
)

func main() {
	cmd, err := commander.Init()
	if err != nil {
		log.Panic(err)
	}
	handlers.Register(cmd)
	cmd.Run()
}
