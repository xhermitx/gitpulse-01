package server

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xhermitx/gitpulse-01/profiler/types"
)

const (
	QUEUE__JOB_STATUS = "JOB_STATUS_QUEUE"
)

type Server struct {
	Store types.CandidateStore
	Git   types.GitService
	Queue types.Queue
	Cache types.Cache
}

func NewServer(s types.CandidateStore, g types.GitService, q types.Queue, c types.Cache) Server {
	return Server{
		Store: s,
		Git:   g,
		Queue: q,
		Cache: c,
	}
}

func (s Server) Run() error {
	msgs, err := s.Queue.Subscribe(QUEUE__JOB_STATUS)
	if err != nil {
		return err
	}

	go s.handleQueueData(msgs)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-msgs
	return nil
}

func (s Server) handleQueueData(msgs <-chan amqp.Delivery) {
	if msgs == nil {
		log.Println("empty body from queue")
	}

	for d := range msgs {
		_ = d
	}
}
