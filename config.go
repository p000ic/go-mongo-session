package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Config mongodb configuration parameters
type Config struct {
	URL           string
	Database      string
	Collection    string
	AuthMechanism string
	Username      string
	Password      string
	AuthSource    string
	Auth          bool
	ClientOptions *options.ClientOptions
}

// NewConfig create mongodb configuration
func NewConfig(url, database, collection, username, password, authSource string) *Config {
	maxConnIdleTime := time.Duration(1000) * time.Millisecond
	return &Config{
		URL:           url,
		Database:      database,
		Collection:    collection,
		AuthMechanism: "SCRAM-SHA-1",
		Username:      username,
		Password:      password,
		AuthSource:    authSource,
		Auth:          true,
		ClientOptions: &options.ClientOptions{
			MaxConnIdleTime: &maxConnIdleTime,
		},
	}
}
