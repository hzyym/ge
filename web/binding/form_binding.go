package binding

import (
	"net/http"
	"reflect"
)

type formBinding struct {
}

func (formBinding) Bind(req *http.Request, obj interface{}) error {
	return mapForm(req, reflect.ValueOf(obj), "form")
}

func mapForm(req *http.Request, obj reflect.Value, tag string) error {
	value := obj
	kind := obj.Kind()
	if kind == reflect.Pointer {
		return mapForm(req, value.Elem(), "form")
	}
	if kind == reflect.Struct {
		//指针方式
		Type := value.Type()
		for i := 0; i < Type.NumField(); i++ {
			field := Type.Field(i)

			v, ok := field.Tag.Lookup(tag)
			if !ok {
				continue
			}
			err := setFieldValue(value.Field(i), req.Form.Get(v))
			if err != nil {
				return err
			}
		}
	} else if kind == reflect.Map {
		mapData := reflect.MakeMap(value.Type())

		for k, v := range req.Form {
			mapData.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v[len(v)-1]))
		}
		value.Set(mapData)
	} else if kind == reflect.Slice {
		sliceValue := reflect.MakeSlice(value.Type(), len(req.Form), len(req.Form))
		var i int
		for _, val := range req.Form {
			sliceValue.Index(i).SetString(val[len(val)-1])
			i++
		}
		value.Set(sliceValue)
	}
	return nil
}
func setFieldValue(value reflect.Value, data string) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(data, 0, value)
	case reflect.String:
		value.SetString(data)
		return nil
	case reflect.Int8:
		return setIntField(data, 8, value)
	case reflect.Int16:
		return setIntField(data, 16, value)
	case reflect.Int32:
		return setIntField(data, 32, value)
	case reflect.Int64:
		return setIntField(data, 64, value)
	case reflect.Bool:
		return setBoolField(data, value)
	case reflect.Float32:
		return setFloatField(data, 32, value)
	case reflect.Float64:
		return setFloatField(data, 64, value)
	case reflect.Uint:
		return setUintField(data, 0, value)
	case reflect.Uint8:
		return setUintField(data, 8, value)
	case reflect.Uint16:
		return setUintField(data, 16, value)
	case reflect.Uint32:
		return setUintField(data, 32, value)
	case reflect.Uint64:
		return setUintField(data, 64, value)
	default:

	}

	return nil
}
