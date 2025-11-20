package db

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/yi-cloud/rest-server/pkg/logs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

type DBManager struct {
	db        *gorm.DB
	dsn       string
	tabModels []interface{}
}

var (
	_db *DBManager
)

func DBInstance(dsn ...string) *DBManager {
	if _db == nil {
		_db = &DBManager{}
	}

	if len(dsn) > 0 {
		_db.dsn = dsn[0]
		logs.Logger.Info(_db.dsn)
	}
	return _db
}

func GormDB() *gorm.DB {
	return DBInstance().db
}

func GetDBLogLevel() logger.LogLevel {
	level := strings.ToLower(viper.GetString("database.loglevel"))
	switch level {
	case "silent":
		return logger.Silent
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Error
	}
}

func (m *DBManager) GormDB() *gorm.DB {
	return m.db
}

func (m *DBManager) Open(t string) *DBManager {
	var dialector gorm.Dialector
	switch t {
	case "mysql":
		dialector = mysql.Open(m.dsn)
	case "sqlite":
		dialector = sqlite.Open(m.dsn)
	case "postgres":
		dialector = postgres.Open(m.dsn)
	default:
		panic("not support database type: " + t)
	}

	dbLogger := logger.New(logs.Logger, logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  GetDBLogLevel(),
		IgnoreRecordNotFoundError: false,
		Colorful:                  false,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	m.db = db
	return m
}

func (m *DBManager) AutoMigrate(dst ...interface{}) {
	err := m.db.AutoMigrate(dst)
	if err != nil {
		panic("auto migrate table failed")
	}
}

func (m *DBManager) AddTabModel(tabModels ...interface{}) {
	if len(tabModels) > 0 {
		m.tabModels = append(m.tabModels, tabModels...)
	}
}

func (m *DBManager) AutoMigrateAll() {
	var oldVersion Version
	if need, err := NewVersionRepository(m.db).IsAutoMigrate(&oldVersion); !need || err != nil {
		return
	}
	err := m.db.AutoMigrate(m.tabModels...)
	if err != nil {
		if oldVersion.Version != "" {
			m.db.Model(&oldVersion).Update("version", oldVersion.Version)
		}
		panic(fmt.Sprintf("auto migrate table failed: %v", err))
	}
}

func getSSLCfg(mode string) string {
	sslcfg := ""
	if mode == "verify-ca" || mode == "verify-full" {
		sslrootcert := viper.GetString("database.sslrootcert")
		if sslrootcert != "" {
			sslcfg += fmt.Sprintf("sslrootcert=%s ", sslrootcert)
		}

		sslcert := viper.GetString("database.sslcert")
		if sslcert != "" {
			sslcfg += fmt.Sprintf("sslcert=%s ", sslcert)
		}

		sslkey := viper.GetString("database.sslkey")
		if sslkey != "" {
			sslcfg += fmt.Sprintf("sslkey=%s ", sslkey)
		}
	}
	return sslcfg
}

func GetDsnBaseInfo() (string, string, string, string, int, string) {
	user := viper.GetString("database.user")
	if user == "" {
		user = "root"
	}

	password := viper.GetString("database.password")
	db := viper.GetString("database.db")
	if db == "" {
		db = "rest-server"
	}

	host := viper.GetString("database.host")
	if host == "" {
		host = "localhost"
	}

	charset := viper.GetString("database.charset")
	if charset == "" {
		charset = "utf8mb4"
	}

	port := viper.GetInt("database.port")
	return user, password, db, host, port, charset
}

// GetDsnFromConfig
// pg_dsn := "host=host user=user password=password dbname=dbname port=port sslmode=disable TimeZone=Asia/Shanghai"

func GetDsnForPostgres() string {
	user, password, db, host, port, charset := GetDsnBaseInfo()
	if port <= 0 {
		port = 5432
	}

	sslmode := viper.GetString("database.sslmode")
	if sslmode == "" {
		sslmode = "disable"
	}

	timezone := viper.GetString("database.timezone")
	if timezone == "" {
		timezone = "Asia/Shanghai"
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d client_encoding=%s sslmode=%s TimeZone=%s %s",
		user, password, host, db, port, charset, sslmode, timezone, getSSLCfg(sslmode))
}

// dsn := "root:123@tcp(127.0.0.1:3306)/rest-server?charset=utf8mb4&parseTime=True&loc=Local"

func GetDsnForMysql() string {
	user, password, db, host, port, charset := GetDsnBaseInfo()
	if port <= 0 {
		port = 3306
	}

	parseTime := "False"
	if viper.GetBool("database.parsetime") {
		parseTime = "True"
	}

	local := ""
	if viper.GetBool("database.local") {
		local = "Local"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s",
		user, password, host, port, db, charset, parseTime)
	if local != "" {
		dsn += "&loc=" + local
	}
	return dsn
}

// dsn := "rest-server.db"
// only use for test

func GetDsnForSqlite() string {
	db := viper.GetString("database.db")
	if db == "" {
		db = "rest-server.db"
	}

	cached := viper.GetString("database.cached")
	if cached == "" {
		cached = "shared"
	}
	return fmt.Sprintf("file:%s?mode=rw&cached=%s", db, cached)
}

func GetDsnFromConfig(t string) string {
	switch t {
	case "mysql":
		return GetDsnForMysql()
	case "sqlite":
		return GetDsnForSqlite()
	case "postgres":
		return GetDsnForPostgres()
	default:
		panic(fmt.Sprintf("Invalid db type: %s", t))
	}
}

func IsRecordNotFound(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrRecordNotFound)
}
