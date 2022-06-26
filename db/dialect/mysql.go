package dialect

import (
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

type Mysql struct {
	name string
	dns  string
}

func NewMysql(dns string) *Mysql {
	return &Mysql{
		name: "mysql",
		dns:  dns,
	}
}
func (m *Mysql) DriveName() string {
	return m.name
}

func (m *Mysql) DriveData() string {
	return m.dns
}

func (m *Mysql) FieldType(field reflect.Value) string {
	switch field.Kind() {
	case reflect.String:
		return "varchar(255)"
	case reflect.Int, reflect.Int32, reflect.Int8, reflect.Int16, reflect.Int64:
		return "int(11)"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Bool:
		return "bool"
	}
	return ""
}
