package main

import (
	env "github.com/caarlos0/env/v6"
)

// Config structured set of environment based configurations for cottontail
type Config struct {
	RabbitAPIHost       string `env:"RABBIT_API_HOST" envDefault:"http://104.197.50.160:15672"`
	RabbitUsername      string `env:"RABBIT_USERNAME" envDefault:"user"`
	RabbitPassword      string `env:"RABBIT_PASSWORD" envDefault:"GMXL5qKD"`
	MonitorInterval     int    `env:"MONITOR_INTERVAL" envDefault:"250"`
	SlackWebhookURL     string `env:"SLACK_POST_URL" envDefault:"https://hooks.slack.com/services/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"`
	ReadyAlertThreshold int    `env:"READY_ALERT_THRESHOLD" envDefault:"100"`
	AlertConfigFilePath string `env:"ALERT_CONFIG_FILE_PATH" envDefault:"./example_alert_config.json"`
}

// GetConfig parses and returns a new instance of Config
func GetConfig() (Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return cfg, err
}
