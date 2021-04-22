package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

// SlackClient is the global reference to the slack client
var SlackClient *slackClient

// InitSlackClient initiates the slack client
func InitSlackClient() {
	SlackClient = &slackClient{
		WebHookUrl: viper.GetString("slack.feedback_webhook_url"),
	}
}

// SendMessage sends a message using the webhook provided in the config
// the webhook provided should be tied to a bot and a channel
func (sc *slackClient) SendMessage(msg string) error {
	sm := SlackMessage{
		Text: msg,
	}
	return sc.sendHttpRequest(sm)
}

// DefaultSlackTimeout is the default timeout for sending a message
const DefaultSlackTimeout = 5 * time.Second

// slackClient is the struct to represent the the slack client
type slackClient struct {
	WebHookUrl string
	TimeOut    time.Duration
}

// SlackMessage is the struct to represent the slack message
type SlackMessage struct {
	Text string `json:"text,omitempty"`
}

// sendHttpRequest posts the slack message to the webhook url, thus
// sending the message
func (sc slackClient) sendHttpRequest(message SlackMessage) error {
	slackBody, _ := json.Marshal(message)
	req, err := http.NewRequest(http.MethodPost, sc.WebHookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	if sc.TimeOut == 0 {
		sc.TimeOut = DefaultSlackTimeout
	}
	client := &http.Client{Timeout: sc.TimeOut}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	if buf.String() != "ok" {
		return fmt.Errorf("non-ok returned from slack url: %s", buf.String())
	}
	return nil
}
