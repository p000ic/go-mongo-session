***REMOVED***

***REMOVED***
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
***REMOVED***

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
	ClientOptions *options.ClientOptions
***REMOVED***

// NewConfig create mongodb configuration
func NewConfig(url, source, collection, username, password, authSource string***REMOVED*** *Config {
	maxConnIdleTime := time.Duration(1000***REMOVED*** * time.Millisecond
	return &Config{
		URL:           url,
		Source:        source,
		Collection:    collection,
		AuthMechanism: "SCRAM-SHA-1",
		Username:      username,
		Password:      password,
		AuthSource:    authSource,
		Auth:          true,
		ClientOptions: &options.ClientOptions{
			MaxConnIdleTime: &maxConnIdleTime,
		***REMOVED***,
	***REMOVED***
***REMOVED***
