package mongo

// Config mongodb configuration parameters
type Config struct {
	URL           string
	Source        string
	Collection    string
	AuthMechanism string
	Username      string
	Password      string
	AuthSource    string
	Auth          bool
}

// NewConfig create mongodb configuration
func NewConfig(url, source, collection, username, password, authSource string) *Config {
	return &Config{
		URL:           url,
		Source:        source,
		Collection:    collection,
		AuthMechanism: "SCRAM-SHA-1",
		Username:      username,
		Password:      password,
		AuthSource:    authSource,
		Auth:          false,
	}
}
