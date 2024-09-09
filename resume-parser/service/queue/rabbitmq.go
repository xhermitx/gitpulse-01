package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xhermitx/gitpulse-01/resume-parser/config"
	"github.com/xhermitx/gitpulse-01/resume-parser/types"
)

type RabbitMQ struct {
	Channel *amqp.Channel
}

func NewRabbitMQClient(ch *amqp.Channel) *RabbitMQ {
	return &RabbitMQ{
		Channel: ch,
	}
}

// FUNCTION TO PUSH CANDIDATE DATA TO THE MESSAGE QUEUE
func (mq *RabbitMQ) Publish(queueName string, data any) error {

	q, err := mq.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(data)

	failOnError(err, fmt.Sprintf("Failed to Parse Status for: %s", data.(types.StatusQueue).JobId))

	if err = mq.Channel.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		}); err != nil {
		return err
	}

	log.Printf("\n[x] Sent status for job %s", data.(types.StatusQueue).JobId)

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func Connect(maxRetries int64) (*amqp.Channel, error) {
	var (
		counts     int64
		backOff    = 1 * time.Second
		connection *amqp.Connection
	)

	log.Println(os.Getenv("RABBITMQ"))

	// Retry with exponential timeout
	for {
		c, err := amqp.Dial(config.Envs.RMQAddr)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > maxRetries {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(2, float64(counts))) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	return channel, nil
}
