package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DBJson []map[string]interface{}

func (j *DBJson) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

func (j DBJson) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j DBJson) GetString(index uint, key string) string {
	if v, ok := j[index][key].(string); ok {
		return v
	}
	return ""
}

func (j DBJson) GetInt(index uint, key string) int {
	if v, ok := j[index][key].(float64); ok {
		return int(v)
	}
	return 0
}

// GormDataType 指定GORM数据类型
func (DBJson) GormDataType() string {
	return "json"
}

// GormDBDataType 指定数据库数据类型
func (DBJson) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// returns different database type based on driver name
	switch db.Dialector.Name() {
	case "mysql":
		return "json"
	case "sqlite":
		return "text"
	case "postgres":
		return "jsonb"
	}
	return ""
}
