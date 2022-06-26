package db

import (
	"fmt"
	"gz/db/dialect"
	"testing"
)

type TestInsert struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func TestDb(t *testing.T) {

	db, err := Open(dialect.NewMysql("root:123456789@tcp(127.0.0.1:3306)/kz"))
	if err != nil {
		t.Error(err)
		return
	}

	var tmp TestInsert
	db.Table("test").Where("name", "opz0").First(&tmp)
	fmt.Println(tmp)

}
