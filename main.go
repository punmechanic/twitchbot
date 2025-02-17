package main

import (
	"context"
	"log"

	"example.com/twitchbot/pkg/twitchbot"
)

func main() {
	ctx := context.Background()
	err := twitchbot.Run(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
