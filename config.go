***REMOVED***

// Config mongodb configuration parameters
type Config struct {
	URL           string
	DB            string
	Collection    string
	AuthMechanism string
	Username      string
	Password      string
	Source        string
	Auth          bool
***REMOVED***

// NewConfig create mongodb configuration
func NewConfig(url, db, collection, username, password, authsource string***REMOVED*** *Config {
	return &Config{
		URL:           url,
		DB:            db,
		Collection:    collection,
		AuthMechanism: "SCRAM-SHA-1",
		Username:      username,
		Password:      password,
		Source:        authsource,
		Auth:          false,
	***REMOVED***
***REMOVED***
