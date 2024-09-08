package main

import (
	"fmt"
	"log"

	"github.com/xhermitx/gitpulse-01/resume-parser/cmd/api"
	"github.com/xhermitx/gitpulse-01/resume-parser/config"
)

func main() {
	server := api.NewAPIServer(fmt.Sprintf("%s:%s", config.Envs.PublicHost, config.Envs.Port))

	log.Fatal(server.Run())

}
