package queue

import (
	"fmt"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ encapsulates a connection channel to RabbitMQ.
type RabbitMQ struct {
	Channel *amqp.Channel
}

// NewRabbitMQClient initializes a new RabbitMQ client with an existing channel.
func NewRabbitMQClient(ch *amqp.Channel) *RabbitMQ {
	return &RabbitMQ{Channel: ch}
}

func (rmq RabbitMQ) Subscribe(queueName string) (<-chan amqp.Delivery, error) {
	q, err := rmq.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	msgs, err := rmq.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

// Connect establishes a RabbitMQ connection and channel with retry logic.
func RMQConnect(maxRetries int64, addr string) (*amqp.Channel, error) {
	var (
		counts     int64             // retry attempt counter
		backOff    = 1 * time.Second // initial backoff duration
		connection *amqp.Connection  // RabbitMQ connection
	)

	// Retry connecting to RabbitMQ with exponential backoff.
	for {
		conn, err := amqp.Dial(addr)
		if err != nil {
			// If connection fails, increment the retry counter and apply backoff.
			log.Println("RabbitMQ not ready, retrying...")
			counts++
			if counts > maxRetries {
				return nil, fmt.Errorf("reached max retries: %w", err)
			}
			// Exponential backoff
			backOff = time.Duration(math.Pow(2, float64(counts))) * time.Second
			log.Printf("backing off for %v...\n", backOff)
			time.Sleep(backOff)
		} else {
			// Successful connection
			log.Println("Connected to RabbitMQ!")
			connection = conn
			break
		}
	}

	// Open a channel on the RabbitMQ connection.
	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return channel, nil
}
