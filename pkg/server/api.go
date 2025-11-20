package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yi-cloud/rest-server/pkg/config"
	"github.com/yi-cloud/rest-server/pkg/logs"
	"github.com/yi-cloud/rest-server/pkg/middleware"
)

type MiddlewareFunc func(*ApiServer)

type Middleware struct {
	name     string
	handlers []MiddlewareFunc
}

type ApiServer struct {
	r           *gin.Engine
	api         *gin.RouterGroup
	version     string
	release     bool
	middlewares []Middleware
}

func New(apiVer ...string) *ApiServer {
	_apiVer := "/api/v1"
	if len(apiVer) > 0 {
		_apiVer = apiVer[0]
	}
	return &ApiServer{version: _apiVer, release: true}
}

func (s *ApiServer) Start() *ApiServer {
	if s.r != nil {
		return s
	}
	gin.SetMode(gin.ReleaseMode)
	s.r = gin.New()
	s.api = s.r.Group(s.version)
	for _, m := range s.middlewares {
		for _, handler := range m.handlers {
			handler(s)
		}
	}
	return s
}

func (s *ApiServer) SetMode(mode string) *ApiServer {
	gin.SetMode(mode)
	if mode != gin.ReleaseMode {
		s.release = false
	}
	return s
}

func (s *ApiServer) RegistryMiddlewares(name string, handlers ...MiddlewareFunc) *ApiServer {
	s.middlewares = append(s.middlewares, Middleware{name: name, handlers: handlers})
	return s
}

func (s *ApiServer) Run() {
	var err error
	domain := viper.GetString("server.ip") + ":" + viper.GetString("server.port")
	logs.Logger.Infof("console server started: %s", domain)
	if config.IsHttpScheme {
		err = s.r.Run(domain)
	} else {
		err = s.r.RunTLS(domain, config.HttpsCrt, config.HttpsKey)
	}
	if err != nil {
		if err != nil {
			panic(fmt.Sprintf("failed to boot server: %s", err.Error()))
		}
	}
}

func (s *ApiServer) GetBaseGroup() *gin.RouterGroup {
	return &s.r.RouterGroup
}

func (s *ApiServer) AddGroup(relativePath string, g *gin.RouterGroup) *gin.RouterGroup {
	if g == nil {
		g = s.api
	}
	return g.Group(relativePath)
}

func (s *ApiServer) AddRoute(g *gin.RouterGroup, method, relativePath string, handlers ...gin.HandlerFunc) {
	if g == nil {
		g = s.api
	}

	switch method {
	case "GET":
		g.GET(relativePath, handlers...)
	case "POST":
		g.POST(relativePath, handlers...)
	case "DELETE":
		g.DELETE(relativePath, handlers...)
	case "PATCH":
		g.PATCH(relativePath, handlers...)
	case "PUT":
		g.PUT(relativePath, handlers...)
	case "OPTIONS":
		g.OPTIONS(relativePath, handlers...)
	case "HEAD":
		g.HEAD(relativePath, handlers...)
	case "Any":
		g.Any(relativePath, handlers...)
	}
}

func PrintBodyMiddleware(s *ApiServer) {
	// for all path registry print body middleware
	if !s.release {
		s.GetBaseGroup().Use(middleware.PrintRawBody())
	}
}

func LogMiddleware(s *ApiServer) {
	// for all path registry log middleware
	g := s.GetBaseGroup()
	g.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output:    logs.AccessLogger.Writer(),
		Formatter: logs.AccessInfo,
	}))
	g.Use(gin.CustomRecovery(logs.RecoveryError))
}

func AuthMiddleware(s *ApiServer) {
	// only for api group registry auth middleware
	middleware.InitRsaKey()
	s.api.Use(middleware.AuthMiddleware())
}
