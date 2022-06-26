package session

import (
	"testing"
)

type TestInsert struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func (TestInsert) TableName() string {
	return "test"
}
func TestSession_Open(t *testing.T) {

}
