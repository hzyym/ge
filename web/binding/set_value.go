package binding

import (
	"reflect"
	"strconv"
)

func setIntField(val string, bit int, field reflect.Value) error {
	v, err := strconv.ParseInt(val, 10, bit)
	if err != nil {
		return err
	}
	field.SetInt(v)
	return nil
}
func setFloatField(val string, bit int, field reflect.Value) error {
	v, err := strconv.ParseFloat(val, bit)
	if err != nil {
		return err
	}
	field.SetFloat(v)
	return err
}
func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	v, err := strconv.ParseBool(val)
	if err != nil {
		return err
	}
	field.SetBool(v)
	return nil
}
func setUintField(val string, bit int, field reflect.Value) error {
	v, err := strconv.ParseUint(val, 10, bit)
	if err != nil {
		return err
	}
	field.SetUint(v)
	return nil
}
