package common

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func ValidateJsonDateType(field reflect.Value) interface{} {
	if field.Type() == reflect.TypeOf(MyTime{}) {
		timeStr := field.Interface().(MyTime).String()
		// 0001-01-01 00:00:00 是 go 中 time.Time 类型的空值
		// 这里返回 Nil 则会被 validator 判定为空值，而无法通过 `binding:"required"` 规则
		if timeStr == "0001-01-01 00:00:00" {
			return nil
		}
		return timeStr
	}
	return nil
}

func AddJsonDateTypeValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterCustomTypeFunc(ValidateJsonDateType, MyTime{})
	}
}

func AddValidatorForServer() {
	AddJsonDateTypeValidator()
}
