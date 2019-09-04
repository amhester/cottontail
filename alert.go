package main

// AlertTrigger string alias for the trigger enum
type AlertTrigger string

// AlertTrigger enum values
const (
	AlertTriggerThreshold = "threshold"
	AlertTriggerDeviation = "deviation"
	AlertTriggerConsumers = "consumers"
)

// Alert contains the contextual information needed when an alert is sent
type Alert struct {
	Trigger AlertTrigger
	AlertConfig
	Queue RabbitQueue
	Stats StatTracker
}
