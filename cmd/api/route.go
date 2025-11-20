package main

import (
	"github.com/yi-cloud/rest-server/api/v1/handlers"
	"github.com/yi-cloud/rest-server/pkg/server"
)

func AddAllRoutes(s *server.ApiServer) {
	s.AddRoute(s.GetBaseGroup(), "POST", "/login", handlers.Login)
	handlers.AddUserRoutes(s)
	handlers.AddUserProfileRoutes(s)
}
