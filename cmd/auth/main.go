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
		SetMode(config.RunMode).
		RegistryMiddlewares("auth", server.LicenseMiddleware, server.AuthMiddleware).
		RegistryMiddlewares("logs", server.LogMiddleware).
		RegistryMiddlewares("printBody", server.PrintBodyMiddleware).
		Start()
	AddAllRoutes(s)
	s.Run()
}
