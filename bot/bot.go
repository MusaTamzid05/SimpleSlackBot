package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
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

func (s *SlackBot) prepareFileParameters(path string) (map[string]io.Reader, error) {

	values := make(map[string]io.Reader)

	fp, err := os.Open(path)

	if err != nil {
		return values, err
	}

	values["file"] = fp
	values["channels"] = strings.NewReader(s.channelName)

	return values, nil
}

func (s *SlackBot) prepareFileData(values map[string]io.Reader) (bytes.Buffer, string, error) {

	var buffer bytes.Buffer
	var err error
	writer := multipart.NewWriter(&buffer)
	var fw io.Writer

	for key, reader := range values {

		fp, ok := reader.(*os.File)

		if ok {
			if fw, err = writer.CreateFormFile(key, fp.Name()); err != nil {
				return buffer, "", err
			}
		} else {
			if fw, err = writer.CreateFormField(key); err != nil {
				return buffer, "", err
			}
		}

		if _, err := io.Copy(fw, reader); err != nil {
			return buffer, "", err
		}

	}

	return buffer, writer.FormDataContentType(), nil
}

func (s *SlackBot) UploadFile(path string) error {
	// https://stackoverflow.com/questions/20205796/post-data-using-the-content-type-multipart-form-data
	values, err := s.prepareFileParameters(path)

	if err != nil {
		return err
	}

	buffer, contentType, err := s.prepareFileData(values)

	req, err := http.NewRequest("POST", "https://slack.com/api/files.upload", &buffer)

	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.botKey))
	req.Header.Set("Content-Type", contentType)
	resBody, err := s.makeRequest(req)

	if err != nil {
		return err
	}

	fmt.Println(resBody)

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
