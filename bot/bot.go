package bot

import (
	"errors"
	"os"
)

type SlackBot struct {
	channelName string
}

func (s *SlackBot) init(environmentName string) error {
	value := os.Getenv(environmentName)

	if value == "" {
		return errors.New("environ variable not found")
	}

	return nil
}

func (s *SlackBot) SentMessage(message string) {

}

func NewSlackBot(environmentName, channelName string) (*SlackBot, error) {
	bot := SlackBot{}
	bot.channelName = channelName

	err := bot.init(environmentName)

	if err != nil {
		return nil, err
	}
	return &bot, nil
}
