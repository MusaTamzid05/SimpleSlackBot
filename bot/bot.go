package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type SlackBot struct {
	channelName string
	timeout     int64
	botKey      string
}

func (s *SlackBot) init(environmentName string) error {
	value := os.Getenv(environmentName)

	if value == "" {
		return errors.New("environ variable not found")
	}

	s.botKey = value

	return nil
}

func (s *SlackBot) SendMessage(message string) error {

	reqBody, err := json.Marshal(map[string]string{
		"channel": fmt.Sprintf("#%s", s.channelName),
		"text":    message,
	})

	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(reqBody))

	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json;charset=utf-8")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.botKey))
	resStr, err := s.makeRequest(request)

	if err != nil {
		return err
	}

	fmt.Println(resStr)
	return nil
}

func (s *SlackBot) makeRequest(request *http.Request) (string, error) {

	timeout := time.Duration(s.timeout)
	client := http.Client{
		Timeout: timeout * time.Second,
	}

	res, err := client.Do(request)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(body), err
}

func NewSlackBot(environmentName, channelName string, timeout int64) (*SlackBot, error) {
	bot := SlackBot{}
	bot.channelName = channelName
	bot.timeout = timeout

	err := bot.init(environmentName)

	if err != nil {
		return nil, err
	}
	return &bot, nil
}
