package main

import (
	"github.com/yi-cloud/rest-server/pkg/config"
	"github.com/yi-cloud/rest-server/pkg/license"
	"github.com/yi-cloud/rest-server/pkg/server"
)

func main() {
	server.Init()
	license.CheckLicense()
	s := server.New().
		RegistryMiddlewares("logs", server.LogMiddleware).
		RegistryMiddlewares("printBody", server.PrintBodyMiddleware).
		RegistryMiddlewares("auth", server.LicenseMiddleware, server.AuthMiddleware).
		SetMode(config.RunMode).
		Start()
	AddAllRoutes(s)
	server.Run(s)
}
