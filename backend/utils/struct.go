package utils

import (
	"errors"
	"reflect"
)

func isZero(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

func IsStructFull(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return errors.New("object is not a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Struct:
			if err := IsStructFull(field.Interface()); err != nil {
				return errors.New("struct is not full")
			}
		default:
			if isZero(field) {
				return errors.New("struct is not full")
			}
		}
	}
	return nil
}
