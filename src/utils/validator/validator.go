package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Validator struct {}
type Rule map[string]string

func (Validator) Validate(data interface{}, rules Rule) (error error) {
	typ := reflect.TypeOf(data)
	val := reflect.ValueOf(data)

	kd := val.Kind()
	if kd != reflect.Struct {
		error = errors.New("expect struct")
		return
	}

	num := val.NumField()
	for index :=0; index < num; index ++ {
		currentTag := typ.Field(index)
		currentVal := val.Field(index)

		if itemRule := rules[currentTag.Name]; len(itemRule) > 0 {
			itemRuleSlice := strings.Split(itemRule, ";")
			for _, rule := range itemRuleSlice {
				switch rule {
				case "required":
					if IsEmpty(currentVal) {
						error = errors.New(fmt.Sprintf("%s 不能为空", currentTag.Name))
						return
					}
				case "email":
					if false == IsCorrectEmail(currentVal) {
						error = errors.New(fmt.Sprintf("%s 邮箱格式不正确", currentTag.Name))
						return
					}

				}
			}
		}
	}
	return
}

func IsEmpty(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.String:
		return val.Len() == 0
	case reflect.Bool:
		return !val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return val.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return val.IsNil()
	}

	return reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
}

func IsCorrectEmail(val reflect.Value) bool {
	return true
}