package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func CheckFieldsForNil(s any) error {
	var serverValue reflect.Value
	if reflect.TypeOf(s).Kind() == reflect.Ptr {
		serverValue = reflect.ValueOf(s).Elem()
	} else {
		serverValue = reflect.ValueOf(s)
	}
	for i := 0; i < serverValue.NumField(); i++ {
		fieldValue := serverValue.Field(i)
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			return errors.New(fmt.Sprint("field ", serverValue.Type().Field(i).Name, " is nil"))
		}
	}
	return nil
}
