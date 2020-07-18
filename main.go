package main

import (
	"fmt"
	"slack_bot/bot"
)

func main() {

	bot, err := bot.NewSlackBot("SLACK_BOT", "testbot", 120)

	if err != nil {
		fmt.Println(err)
	}

	err = bot.SendMessage("Hello from golang")

	if err != nil {
		fmt.Println(err)
	}
}
