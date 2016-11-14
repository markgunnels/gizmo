package amqp

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"

	"github.com/NYTimes/gizmo/pubsub"
	"github.com/NYTimes/gizmo/web"
)

// URL will return the connection URL AMQP
func (c Config) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		a.User,
		a.Pw,
		a.Host,
		a.Port,
	)
}

type AMQPPublisher struct {
	crypt *Encrypter

	cfg  Config
	conn *amqp.Connection
	chn  *amqp.Channel
}

func NewAMQPPublisher(cfg Config) (pubsub.Publisher, error) {
	pub := &AMQPPublisher{}

	var err error
	pub.conn, err = amqp.Dial(cfg.URL())
	if err != nil {
		Log.Error("unable to connect to AMQP: ", err)
		return pub, err
	}

	pub.chn, err = pub.conn.Channel()
	if err != nil {
		Log.Error("unable to get to AMQP channel: ", err)
		return pub, err
	}

	// make sure the topic is there :)
	err = ch.ExchangeDeclare(
		cfg.Topic,      // name
		"topic",        // type
		cfg.Durable,    // durable
		cfg.AutoDelete, // auto-deleted
		false,          // internal
		cfg.NowWait,    // no-wait
		nil,            // arguments
	)

	if err != nil {
		Log.Error("unable to declare AMQP topic: ", err)
	}

	return pub, err
}

// Publish will marshal the proto message and emit it to the AMQP topic.
func (p *AMQPPublisher) Publish(key string, m proto.Message) error {
	mb, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return p.PublishRaw(key, mb)
}

// Publish will emit the byte array to the AMQP topic.
func (p *AMQPPublisher) PublishRaw(key string, m []byte) error {
	msg := amqp.Publishing{
		DeliveryMode: p.cfg.DeliveryMode,
		Timestamp:    time.Now(),
		ContentType:  web.JSONContentType,
		Body:         m,
	}

	return p.chn.Publish("", p.cfg.Topic, AMQPMandatory, AMQPImmediate, msg)
}

// Stop will close the AMQP channel and connection
func (p *AMQPPublisher) Stop() error {
	p.chn.Close()
	return p.conn.Close()
}

// Subscriber is a generic interface to encapsulate how we want our subscribers
// to behave. For now the system will auto stop if it encounters any errors. If
// a user encounters a closed channel, they should check the Err() method to see
// what happened.
type AMQPSubscriber struct {
	cfg  Config
	conn *amqp.Connection
	chn  *amqp.Channel
	stop chan *amqp.Error
	err  Error
}

func NewAMQPSubscriber(cfg Config) (pubsub.Subscriber, error) {
	sub := &AMQPSubscriber{}
	var err error

	sub.stop = make(chan *amqp.Error)

	sub.conn, err = amqp.Dial(cfg.URL())
	if err != nil {
		Log.Error("unable to connect to AMQP: ", err)
		return sub, err
	}

	sub.chn, err = sub.conn.Channel()
	if err != nil {
		Log.Error("unable to get to AMQP channel: ", err)
		return sub, err
	}

	_, err = sub.chn.QueueDeclare(
		cfg.Queue,
		cfg.Durable,
		cfg.AutoDelete,
		cfg.Exclusive,
		cfg.NoWait,
		nil)
	if err != nil {
		Log.Error("unable to declare AMQP queue: ", err)
		return sub, err
	}

	return sub, nil
}

// Start will return a channel of raw messages.
func (s *AMQPSubscriber) Start() <-chan SubscriberMessage {
	c := make(chan SubscriberMessage)

	go func() {
		msgs, err := s.chn.Consume(
			s.cfg.Queue,     // queue
			"",              // consumer
			false,           // auto-ack
			s.cfg.Exclusive, // exclusive
			false,           // no-local
			s.cfg.NoWait,    // no-wait
			nil,             // args
		)

		if err != nil {
			s.err = err
			return
		}

		for {
			select {
			case msg := <-msgs:
				c <- &subMessage{msg}
				continue
			case amqpErr := <-s.stop:
				s.err = amqpErr
				close(c)
				return
			}
		}
	}()

	return c
}

// // Err will contain any errors returned from the consumer connection.
func (s *subscriber) Err() error {
	return s.err
}

// // Stop will initiate a graceful shutdown of the subscriber connection.
func (s *subscriber) Stop() error {
	err := s.chn.Close()
	if err != nil {
		s.err = err
		return err
	}

	err := s.conn.Close()
	if err != nil {
		s.err = err
		return err
	}

	return nil
}

// SubscriberMessage is a struct to encapsulate subscriber messages and provide
// a mechanism for acknowledging messages _after_ they've been processed.
type subMessage struct {
	message *amqp.Delivery
}

func (m *subMessage) Message() []byte {
	return m.subMessage.Body
}

func (m *subMessage) Done() error {
	return m.subMessage.Ack(false)
}

func (m *subMessage) ExtendDoneDeadline(time.Duration) error {
	return nil
}
