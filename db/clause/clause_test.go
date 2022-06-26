package clause

import (
	"testing"
)

func TestClause_Set(t *testing.T) {
	clause := new(Clause)

	clause.Set(INSERT, "user", []string{"name", "age"})
	clause.Set(VALUES, []interface{}{"test", 20})

	sql, vars := clause.Build(INSERT, VALUES)

	t.Log("sql", sql)
	t.Log("vars", vars)
}
