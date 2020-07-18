package bot

import (
	"errors"
	"os"
)

type SlackBot struct {
}

func (s *SlackBot) init(environmentName string) error {
	value := os.Getenv(environmentName)

	if value == "" {
		return errors.New("environ variable not found")
	}

	return nil
}

func NewSlackBot(environmentName string) (*SlackBot, error) {
	bot := SlackBot{}

	err := bot.init(environmentName)

	if err != nil {
		return nil, err
	}
	return &bot, nil
}
