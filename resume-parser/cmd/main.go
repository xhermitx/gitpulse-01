package main

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/xhermitx/gitpulse-01/resume-parser/cmd/api"
	"github.com/xhermitx/gitpulse-01/resume-parser/config"
	"github.com/xhermitx/gitpulse-01/resume-parser/service/cache"
	"github.com/xhermitx/gitpulse-01/resume-parser/service/drive"
	"github.com/xhermitx/my-utils/queue"
)

func main() {

	// Drive Service
	ds, err := drive.NewGoogleService()
	if err != nil {
		log.Fatal(err)
	}
	storage := drive.NewGoogleDrive(ds)

	// Message Broker
	ch, err := queue.RMQConnect(5, config.Envs.RMQAddr)
	if err != nil {
		log.Fatal(err)
	}
	rmq := queue.NewRabbitMQClient(ch)

	// Cache
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Envs.RedisAddr,
		Password: "",
		DB:       0,
	})
	c := cache.NewRedisClient(rdb)

	address := fmt.Sprintf("%s:%s", config.Envs.PublicHost, config.Envs.Port)
	server := api.NewAPIServer(address, storage, rmq, c)
	log.Fatal(server.Run())
}
