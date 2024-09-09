package main

import (
	"fmt"
	"log"

	"github.com/xhermitx/gitpulse-01/resume-parser/cmd/api"
	"github.com/xhermitx/gitpulse-01/resume-parser/config"
	"github.com/xhermitx/gitpulse-01/resume-parser/service/drive"
	"github.com/xhermitx/gitpulse-01/resume-parser/service/queue"
)

func main() {

	// Drive Service
	ds, err := drive.NewGoogleService()
	if err != nil {
		log.Fatal(err)
	}
	storage := drive.NewGoogleDrive(ds)

	// Message Broker
	ch, err := queue.Connect(5)
	if err != nil {
		log.Fatal(err)
	}
	rmq := queue.NewRabbitMQClient(ch)

	address := fmt.Sprintf("%s:%s", config.Envs.PublicHost, config.Envs.Port)

	server := api.NewAPIServer(address, storage, rmq)
	log.Fatal(server.Run())
}
