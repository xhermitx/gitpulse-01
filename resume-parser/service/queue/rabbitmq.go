package queue

import (
	"context"
	"encoding/json"
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

// Publish sends a message to a specified queue in RabbitMQ with a JSON-encoded body.
func (mq *RabbitMQ) Publish(queueName string, data any) error {
	// Declare a queue to ensure the target queue exists.
	q, err := mq.Channel.QueueDeclare(
		queueName, // name of the queue
		false,     // non-durable (messages won't survive a broker restart)
		false,     // auto-delete when unused
		false,     // not exclusive
		false,     // no-wait (do not wait for the server to confirm queue declaration)
		nil,       // no additional arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Create a context with timeout for publishing the message.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Marshal the message data into JSON format.
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Publish the message to the queue.
	if err := mq.Channel.PublishWithContext(
		ctx,
		"",     // exchange (empty string means default exchange)
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json", // MIME type of the message body
			Body:        body,               // message body as a byte slice
		}); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
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
