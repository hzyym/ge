package dialect

import (
	"reflect"
	"testing"
)

func TestMysql_FieldType(t *testing.T) {
	var str float32

	t.Log(NewMysql("test").FieldType(reflect.ValueOf(str)))
}
