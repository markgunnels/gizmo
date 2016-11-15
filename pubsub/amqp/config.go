package amqp

import (
	"github.com/NYTimes/gizmo/config"
)

// Config holds common credentials and config values for
// working with AMQP PubSub.
type Config struct {
	URL   string `envconfig:"AMQP_URL"`
	Topic string `envconfig:"AMQP_TOPIC"`
	Queue string `envconfig:"AMQP_QUEUE"`

	// Delivery mode is the delivery mode used for AMQP publishers
	DeliveryMode uint8 `envconfig:"AMQP_DELIVERY_MODE"`

	// Durable is for the publisher QueueDeclare flags
	Durable bool `envconfig:"AMQP_DURABLE"`
	// AutoDelete is for the publisher QueueDeclare flags
	AutoDelete bool `envconfig:"AMQP_AUTO_DELETE"`
	// Exclusive is for the publisher QueueDeclare flags
	Exclusive bool `envconfig:"AMQP_EXCLUSIVE"`
	// NoWait is for the publisher QueueDeclare flags
	NoWait bool `envconfig:"AMQP_NO_WAIT"`

	// Mandatory is for the 'Publish' flags
	Mandatory bool `envconfig:"AMQP_MANDATORY"`
	// Immediate is for the 'Publish' flags
	Immediate bool `envconfig:"AMQP_IMMEDIATE"`
}

// Attempts to load the AMQP Config struct
// from environment variables. If not populated, nil
// is returned.
func LoadAMQPConfigFromEnv() Config {
	var cfg Config
	config.LoadEnvConfig(&cfg)
	return cfg
}
