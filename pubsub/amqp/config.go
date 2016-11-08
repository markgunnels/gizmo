package amqp

// Config holds common credentials and config values for
// working with AMQP PubSub.
type Config struct {
	Host  string
	Port  int
	User  string
	Pw    string
	Topic string
	Queue string

	// Delivery mode is the delivery mode used for AMQP publishers
	DeliveryMode uint8

	// Durable is for the publisher QueueDeclare flags
	Durable bool
	// AutoDelete is for the publisher QueueDeclare flags
	PAutoDelete bool
	// Exclusive is for the publisher QueueDeclare flags
	Exclusive bool
	// NoWait is for the publisher QueueDeclare flags
	NoWait bool

	// Mandatory is for the 'Publish' flags
	Mandatory bool
	// Immediate is for the 'Publish' flags
	Immediate bool
}

// LoadConfigFromEnv will attempt to load a AMQP config
// from environment variables.
func LoadConfigFromEnv() Config {
	var c Config
	c.LoadEnvConfig(&c)
	return c
}
