package main

import (
	log "github.com/rs/zerolog/log"
)

func main() {
	// Get service configurations from env
	cfg, err := GetConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config from env")
		return
	}

	// Get alert trigger configurations from file
	alertConfigs, err := GetAlertConfig(cfg.AlertConfigFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to retrieve and parse alerts config")
	}

	// Create RabbitMQ Client
	rabbitClient := NewRabbitClient(cfg.RabbitAPIHost, cfg.RabbitUsername, cfg.RabbitPassword)

	// Create Slack Client
	slackClient := NewSlackClient(cfg.SlackWebhookURL)

	// Initialize monitor
	monitor := NewMonitor(&alertConfigs, rabbitClient, slackClient)

	// Start cottontail monitoring rabbitmq
	done := <-monitor.Start()

	// If exited with an error, log it
	if done != nil {
		log.Fatal().Err(done).Msg("Unexpected end of service")
	}
}
