package main

import (
	"fmt"
	"slack_bot/bot"
)

func main() {

	_, err := bot.NewSlackBot("SLACK_BOT")

	if err != nil {
		fmt.Println(err)
	}
}
