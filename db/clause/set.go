package clause

import (
	"fmt"
	"strings"
)

type TypeClause int

const (
	INSERT TypeClause = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDER
)

type setCall func(values ...interface{}) (string, []interface{})

var setCalls = map[TypeClause]setCall{
	INSERT: insert,
	VALUES: values,
	SELECT: select_,
	WHERE:  where,
	ORDER:  order,
	LIMIT:  limit,
}

func setBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}
func insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%s)", tableName, fields), nil
}
func values(values ...interface{}) (string, []interface{}) {
	var bindstr string
	var writeStr strings.Builder
	var vars []interface{}

	writeStr.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindstr == "" {
			bindstr = setBindVars(len(v))
		}
		writeStr.WriteString(fmt.Sprintf("(%s)", bindstr))
		if i+1 != len(values) {
			writeStr.WriteString(",")
		}
		vars = append(vars, v...)
	}
	return writeStr.String(), vars
}
func select_(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %s FROM %s", fields, tableName), nil
}
func where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s=?", desc), vars
}
func order(values ...interface{}) (string, []interface{}) {
	fields := strings.Join(values[0].([]string), ",")
	return fmt.Sprintf("ORDER BY %s", fields), nil
}
func limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}
