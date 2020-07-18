package main

import (
	"fmt"
	"slack_bot/bot"
)

func main() {

	_, err := bot.NewSlackBot("SLACK_BOT", "testbot")

	if err != nil {
		fmt.Println(err)
	}
}
