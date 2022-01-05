package config

import "github.com/plieskovsky/go-grpc-server-shop/internal/server"

// Configuration structure.
type Configuration struct {
	Server Servers
}

// Servers configuration structure.
type Servers struct {
	Grpc server.Config
}
