package server

import (
	"flag"
	"fmt"
	"github.com/yi-cloud/rest-server/common"
	"github.com/yi-cloud/rest-server/pkg/config"
	"github.com/yi-cloud/rest-server/pkg/db"
	"github.com/yi-cloud/rest-server/pkg/license"
	"github.com/yi-cloud/rest-server/pkg/logs"
	"os"
)

var (
	Commit    = "none"
	Version   = "dev"
	BuildDate = "unknown"
)

var (
	cfg  string
	help bool
	h    bool
)

func InitOption() {
	flag.StringVar(&cfg, "config", "/etc/rest-server/config.yaml", "app config file")
	flag.BoolVar(&help, "help", false, "this help command")
	flag.BoolVar(&h, "h", false, "this help command")

	flag.Parse()
	if h || help {
		fmt.Printf("Current server version is %s, git version is %s, build date is %s\n",
			Version, Commit, BuildDate)
		if license.ClusterId != "" {
			fmt.Printf("ClusterID is %s\n", license.ClusterId)
		} else {
			fmt.Printf("ClusterID is none\n")
		}
		fmt.Printf("\n")
		flag.Usage()
		fmt.Printf("\n")
		os.Exit(0)
	}
	flag.Usage = usage
}

func InitLog() {
	logCfg := logs.NewLogConfig()
	logCfg.SetLogConfig(config.CommonLogs)
	logs.InitLog(logCfg)

	logs.Logger.SetLevel(logs.GetLogLevel(config.LogLevel))

	accessCfg := logs.NewLogConfig()
	accessCfg.SetLogConfig(config.AccessLogs)
	logs.InitAccessLog(accessCfg)
}

func Init() {
	InitOption()
	config.LoadConfig(cfg)
	InitLog()
	logs.Logger.Debugf("config path: %s", cfg)
}

func Run(s *ApiServer) {
	if config.UseDB {
		db.DBInstance(db.GetDsnFromConfig(config.DBType)).Open(config.DBType).AutoMigrateAll()
	}
	common.AddValidatorForServer()
	s.Run()
}

func usage() {
	flag.PrintDefaults()
}
