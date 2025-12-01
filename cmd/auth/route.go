package main

import (
	"github.com/yi-cloud/rest-server/api/v1/handlers"
	"github.com/yi-cloud/rest-server/pkg/server"
)

func AddAllRoutes(s *server.ApiServer) {
	s.AddRoute(s.GetBaseGroup(), "POST", "/login", handlers.Login)
	s.AddRoute(s.GetBaseGroup(), "POST", "/ec2tokens", handlers.EC2Tokens)
	handlers.AddUserRoutes(s)
	handlers.AddRoleRoutes(s)
	handlers.AddProjectRoutes(s)
}
