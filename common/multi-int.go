package common

import (
	"database/sql/driver"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"strings"
)

type MultiInt []int

func (s *MultiInt) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("failed to scan multistring field - source is not a string")
	}

	strArr := strings.Split(string(b), ",")
	for _, v := range strArr {
		i, _ := strconv.Atoi(v)
		*s = append(*s, i)
	}
	return nil
}

func (s MultiInt) Value() (driver.Value, error) {
	if s == nil || len(s) == 0 {
		return nil, nil
	}
	strArr := make([]string, len(s))
	for i, v := range s {
		strArr[i] = strconv.Itoa(v) // 将整型转换为字符串
	}
	return strings.Join(strArr, ","), nil
}

// GormDataType 指定GORM数据类型
func (MultiInt) GormDataType() string {
	return "text"
}

// GormDBDataType 指定数据库数据类型
func (MultiInt) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql", "sqlite", "postgres":
		return "text"
	}
	return ""
}
