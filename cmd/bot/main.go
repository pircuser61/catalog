package main

import (
	"go1/internal/commander"
	"go1/internal/handlers"
	"log"
)

func main() {
	cmd, err := commander.Init()
	if err != nil {
		log.Panic(err)
	}
	handlers.Register(cmd)
	cmd.Run()
}
