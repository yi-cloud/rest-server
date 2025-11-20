package common

import (
	"database/sql/driver"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

type MultiString []string

func (s *MultiString) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("failed to scan multistring field - source is not a string")
	}
	*s = strings.Split(string(b), ",")
	return nil
}

func (s MultiString) Value() (driver.Value, error) {
	if s == nil || len(s) == 0 {
		return nil, nil
	}
	return strings.Join(s, ","), nil
}

// GormDataType 指定GORM数据类型
func (MultiString) GormDataType() string {
	return "text"
}

// GormDBDataType 指定数据库数据类型
func (MultiString) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql", "sqlite", "postgres":
		return "text"
	}
	return ""
}
