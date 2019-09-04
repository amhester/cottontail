package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	rabbitmqQueueStatsRoute = "/api/queues"
)

// RabbitClient wraps logic for interacting with RabbitMQ
type RabbitClient struct {
	client      *http.Client
	host        string
	user        string
	password    string
	contentType string
}

// NewRabbitClient returns a new instance of a RabbitClient with the given connection params
func NewRabbitClient(host string, user string, password string) *RabbitClient {
	return &RabbitClient{
		client:      http.DefaultClient,
		host:        host,
		user:        user,
		password:    password,
		contentType: defaultContentType,
	}
}

// GetQueueStats retrieves a list of all queues and their metrics from the configured rabbitmq host
func (client *RabbitClient) GetQueueStats() ([]*RabbitQueue, error) {
	url := fmt.Sprintf("%s%s", client.host, rabbitmqQueueStatsRoute)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", client.contentType)
	req.Header.Add("Accept", client.contentType)
	req.Header.Add("User-Agent", defaultUserAgent)
	req.SetBasicAuth(client.user, client.password)
	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	var queues []*RabbitQueue
	if err := json.NewDecoder(res.Body).Decode(&queues); err != nil {
		return nil, err
	}
	return queues, nil
}
