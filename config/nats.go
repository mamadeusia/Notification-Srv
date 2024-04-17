package config

import (
	"strings"
	"time"
)

type Nats struct {
	URL          string
	Notification Notification
	Nkey         string
}

type Notification struct {
	Prefix string
	Retry  Retry
}

type Retry struct {
	Unit      string
	Durations string
}

type Max struct {
	Tx int
}

func NatsNotificationRetryUnitString() string {
	switch cfg.Nats.Notification.Retry.Unit {
	case "Second":
		return "Second"
	case "Minute":
		return "Minute"
	case "Hour":
		return "Hour"
	default:
		return "Second"
	}
}

func NatsCallbackRetryUnit() time.Duration {
	switch cfg.Nats.Notification.Retry.Unit {
	case "Second":
		return time.Second
	case "Minute":
		return time.Minute
	case "Hour":
		return time.Hour
	default:
		return time.Second
	}
}

func NatsNotificationRetryDurations() []string {
	if cfg.Nats.Notification.Retry.Durations == "" {
		return []string{"1", "5", "7", "12", "20"}
	}
	return strings.Split(cfg.Nats.Notification.Retry.Durations, ",")
}

func NatsNotificationPrefix() string {
	return cfg.Nats.Notification.Prefix
}

func NatsNkey() string {
	return cfg.Nats.Nkey
}
func NatsURL() string {
	return cfg.Nats.URL
}
