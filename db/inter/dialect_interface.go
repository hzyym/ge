package inter

import "reflect"

type Drive interface {
	DriveName() string
	DriveData() string
	FieldType(field reflect.Value) string
}
