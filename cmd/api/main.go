package main

import (
	"github.com/yi-cloud/rest-server/pkg/config"
	"github.com/yi-cloud/rest-server/pkg/server"
)

func main() {
	server.Init()
	s := server.New().
		RegistryMiddlewares("auth", server.AuthMiddleware).
		RegistryMiddlewares("logs", server.LogMiddleware).
		RegistryMiddlewares("printBody", server.PrintBodyMiddleware).
		SetMode(config.RunMode).
		Start()
	AddAllRoutes(s)
	server.Run(s)
}
