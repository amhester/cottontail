package main

import (
	"time"
)

// Monitor runs periodic checks on rabbit's queue stats endpoint and sends alerts when triggered.
type Monitor struct {
	done         chan error
	cancel       chan struct{}
	alertConfig  *QueueConfigs
	rabbitClient *RabbitClient
	slackClient  *SlackClient
	stats        map[string]StatTracker
}

// NewMonitor returns a new instance of a Monitor
func NewMonitor(alertConfig *QueueConfigs, rabbitClient *RabbitClient, slackClient *SlackClient) *Monitor {
	return &Monitor{
		alertConfig:  alertConfig,
		rabbitClient: rabbitClient,
		slackClient:  slackClient,
		done:         make(chan error, 1),
		cancel:       make(chan struct{}, 1),
		stats:        map[string]StatTracker{},
	}
}

// Start starts an asyncronous loop on a timed interval, running all checks and alerts. Returns a channel for receiving process errors
func (monitor *Monitor) Start() <-chan error {
	go func(interval time.Duration, cancel <-chan struct{}, m *Monitor) {
		if err := m.process(time.Now()); err != nil {
			m.done <- err
			return
		}
		for {
			select {
			case <-cancel:
				m.done <- nil
				return
			case now := <-time.After(interval):
				err := m.process(now)
				if err != nil {
					m.done <- err
					return
				}
			}
		}
	}(time.Minute*5, monitor.cancel, monitor)
	return monitor.done
}

// Stop cancels the process loop, stopping all periodic checks and alerts
func (monitor *Monitor) Stop() {
	monitor.cancel <- struct{}{}
}

// Close closes all of the monitor's channels
func (monitor *Monitor) Close() {
	close(monitor.cancel)
	close(monitor.done)
}

func (monitor *Monitor) process(now time.Time) error {
	queues, err := monitor.rabbitClient.GetQueueStats()
	if err != nil {
		return err
	}
	for _, queue := range queues {
		// If queue hasn't come in yet, initialize its stats object
		if _, ok := monitor.stats[queue.Name]; !ok {
			monitor.stats[queue.Name] = NewStatTracker(144)
		}
		// Update the stats for the queue
		monitor.stats[queue.Name].Update([]int64{queue.MessagesReady})
		// Check if we should alert for queue, send if yes
		if err := monitor.checkAndSendAlert(queue); err != nil {
			return err
		}
	}
	return nil
}

func (monitor *Monitor) checkAndSendAlert(queue *RabbitQueue) error {
	// Get the alert config for this queue, or default
	var conf *QueueConfig
	if c, ok := monitor.alertConfig.Queues[queue.Name]; ok {
		conf = c
	}
	if conf == nil {
		conf = monitor.alertConfig.Default
	}

	stats := monitor.stats[queue.Name]

	// Check if queue ready msg count exceeds threshold
	if conf.Threshold != nil && conf.Threshold.Enabled && queue.MessagesReady > int64(conf.Threshold.Value) {
		monitor.slackClient.SendAlert(Alert{
			Trigger:     AlertTriggerThreshold,
			AlertConfig: *conf.Threshold,
			Queue:       *queue,
			Stats:       stats,
		})
	}

	// Check if queue ready msg count exceeds 3 deviations above norm
	if conf.Deviation != nil && conf.Deviation.Enabled && (float64(queue.MessagesReady)-stats.Mean)/stats.StdDev > 3 {
		monitor.slackClient.SendAlert(Alert{
			Trigger:     AlertTriggerDeviation,
			AlertConfig: *conf.Deviation,
			Queue:       *queue,
			Stats:       stats,
		})
	}

	// Check if queue has enough consumers, according to config
	if conf.Consumer != nil && conf.Consumer.Enabled && queue.Consumers <= conf.Consumer.Value {
		monitor.slackClient.SendAlert(Alert{
			Trigger:     AlertTriggerConsumers,
			AlertConfig: *conf.Consumer,
			Queue:       *queue,
			Stats:       stats,
		})
	}

	return nil
}
