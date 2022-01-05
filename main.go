package main

import (
	"github.com/plieskovsky/go-grpc-server-shop/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
