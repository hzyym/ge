package clause

import "strings"

type Clause struct {
	sql     map[TypeClause]string
	sqlVars map[TypeClause][]interface{}

	buildType []TypeClause
}

func (c *Clause) Set(types TypeClause, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[TypeClause]string)
		c.sqlVars = make(map[TypeClause][]interface{})
	}

	sql, value := setCalls[types](vars...)
	c.sql[types] = sql
	c.sqlVars[types] = value
}
func (c *Clause) Build(types ...TypeClause) (string, []interface{}) {
	var arr []TypeClause
	if len(types) == 0 {
		arr = c.buildType
	} else {
		arr = types
	}
	var sqls []string
	var vars []interface{}
	for _, clause := range arr {
		if sql, ok := c.sql[clause]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[clause]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
func (c *Clause) Add(types ...TypeClause) {
	c.buildType = append(c.buildType, types...)
}
