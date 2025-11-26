package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var RunMode string
var SysRoot string
var IsHttpScheme = true
var DBType = "mysql"
var HttpsCrt = "/etc/rest-server/certs/server.crt"
var HttpsKey = "/etc/rest-server/certs/server.key"
var CommonLogs map[string]any
var AccessLogs map[string]any
var LogLevel = ""
var TokenLifeTime int64

type AppFilePathType uint

func initScheme() {
	scheme := viper.GetString("server.scheme")
	switch scheme {
	case "https":
		cert := viper.GetString("server.sslcert")
		if cert != "" {
			HttpsCrt = cert
		}
		key := viper.GetString("server.sslkey")
		if key != "" {
			HttpsKey = key
		}
		IsHttpScheme = false
	case "http", "":
	default:
		panic(fmt.Sprintf("Invalid http scheme: %s", scheme))
	}
}

func initDBType() {
	dbtype := viper.GetString("database.type")
	switch dbtype {
	case "mysql", "postgres", "sqlite":
		DBType = dbtype
	case "":
	default:
		panic(fmt.Sprintf("Invalid database type: %s", dbtype))
	}
}

func initLogsCfg() {
	CommonLogs = viper.GetStringMap("logs.common")
	AccessLogs = viper.GetStringMap("logs.access")
	LogLevel = viper.GetString("logs.level")
	if LogLevel == "" {
		LogLevel = "info"
	}
}

func initConfig() {
	root := viper.GetString("server.root")
	if root != "" {
		root += "/"
	}
	SysRoot = root
	RunMode = viper.GetString("server.runmode")
	TokenLifeTime = viper.GetInt64("auth.lifetime")
	if TokenLifeTime <= 0 {
		TokenLifeTime = 900
	}

	initScheme()
	initDBType()
	initLogsCfg()
}

func LoadConfig(cfg string, t ...string) {
	path := filepath.Dir(cfg)
	filename := filepath.Base(cfg)

	if path == "" {
		viper.AddConfigPath(".")
		viper.AddConfigPath("etc/rest-server/")
		viper.AddConfigPath("/etc/rest-server/")
	} else {
		viper.AddConfigPath(path)
	}

	_t := "yaml"
	names := strings.Split(filename, ".")
	namesLen := len(names)
	if namesLen > 1 {
		_t = names[namesLen-1]
	}
	if len(t) > 0 {
		_t = t[0]
	}
	viper.SetConfigType(_t)
	viper.SetConfigName(strings.Join(names[:namesLen-1], "."))
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("read config error: %v", err))
	}

	initConfig()
}
