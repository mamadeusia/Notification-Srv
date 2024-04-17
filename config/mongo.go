package config

import "time"

type Mongo struct {
	URL            string
	RequestTimeOut time.Duration
}

func MongoURL() string {
	return cfg.Mongo.URL
}
