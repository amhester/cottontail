package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	defaultContentType = "application/json"
	defaultUserAgent   = "cottontail"
	msgThresholdAlert  = "*MSG THRESHOLD ALERT*\\n*%s* currently has *%d* messages in ready, exceeding configured threshold of *%d*.\\n\\n*QUEUE STATS*\\n```%s```"
	msgDeviationAlert  = "*MSG DEVIATION ALERT*\\n*%s* currently has *%d* messages in ready, exceeding %d deviation above the norm of *%f*.\\n\\n*QUEUE STATS*\\n```%s```"
	msgConsumerAlert   = "*MSG CONSUMER ALERT*\\n*%s* currently has *%d* consumers, which is fewer than the configured number of *%d*.\\n\\n*QUEUE STATS*\\n```%s```"
)

// SlackClient wraps logic for interacting with slack
type SlackClient struct {
	client      *http.Client
	contentType string
	webhookURL  string
	userAgent   string
}

// NewSlackClient returns a new instance of a slack client using the given webhook URL
func NewSlackClient(webhookURL string) *SlackClient {
	return &SlackClient{
		client:      http.DefaultClient,
		contentType: defaultContentType,
		webhookURL:  webhookURL,
		userAgent:   defaultUserAgent,
	}
}

// SendAlert sends the given alert as a formatted message to the client's webhook URL
func (client *SlackClient) SendAlert(alert Alert) error {
	alertMsg := createSlackAlertMessage(alert)
	return client.sendSlackMessage(alertMsg)
}

// SlackMessage structured definition for a block post to a slack webhook
type SlackMessage struct {
	Type string           `json:"type"`
	Text SlackMessageText `json:"text"`
}

// SlackMessageText the text portion that goes inside a slack message block section
type SlackMessageText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func createSlackAlertMessage(alert Alert) SlackMessage {
	var content string
	switch alert.Trigger {
	case AlertTriggerThreshold:
		content = fmt.Sprintf(msgThresholdAlert, alert.Queue.Name, alert.Queue.MessagesReady, alert.AlertConfig.Value, alert.Queue)
	case AlertTriggerConsumers:
		content = fmt.Sprintf(msgConsumerAlert, alert.Queue.Name, alert.Queue.Consumers, alert.AlertConfig.Value, alert.Queue)
	case AlertTriggerDeviation:
		content = fmt.Sprintf(msgDeviationAlert, alert.Queue.Name, alert.Queue.MessagesReady, alert.AlertConfig.Value, alert.Stats.Mean, alert.Queue)
	}

	return SlackMessage{
		Type: "section",
		Text: SlackMessageText{
			Type: "mrkdwn",
			Text: content,
		},
	}
}

func (client *SlackClient) sendSlackMessage(msg SlackMessage) error {
	encodedBody := new(bytes.Buffer)
	if err := json.NewEncoder(encodedBody).Encode(&msg); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, client.webhookURL, encodedBody)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", client.contentType)
	req.Header.Add("Accept", client.contentType)
	req.Header.Add("User-Agent", defaultUserAgent)
	res, err := client.client.Do(req)
	if err != nil {
		return err
	}
	if err := checkResponse(res); err != nil {
		return err
	}
	return nil
}
