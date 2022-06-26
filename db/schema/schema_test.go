package schema

import (
	"gz/db/dialect"
	"testing"
)

type TestFieldSchema struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
	Des  string `db:"des"`
}

func TestSchema_Parse(t *testing.T) {
	drive := dialect.NewMysql("test")
	schema := Parse(&TestFieldSchema{}, drive)
	t.Log(schema)
}
