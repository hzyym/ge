package session

import (
	"database/sql"
	"gz/db/clause"
	"gz/db/inter"
	"gz/db/schema"
	"reflect"
)

type Session struct {
	db         *sql.DB
	drive      inter.Drive
	clause     *clause.Clause
	paresTable *schema.Schema
	err        error

	tableName string

	buildType []clause.TypeClause
}

func (s *Session) Open(drive inter.Drive) error {
	var err error
	s.db, err = sql.Open(drive.DriveName(), drive.DriveData())
	if err != nil {
		return err
	}

	if err = s.db.Ping(); err != nil {
		return err
	}
	s.drive = drive
	if s.clause == nil {
		s.clause = &clause.Clause{}
	}
	return nil
}
func (s *Session) Insert(vars interface{}) (sql.Result, error) {
	var record []interface{}
	table := s.Model(vars).paresTable
	s.clause.Set(clause.INSERT, table.Name, table.FieldsName)
	record = table.RecordValues(vars)
	s.clause.Set(clause.VALUES, record...)
	sqls, values := s.clause.Build(clause.INSERT, clause.VALUES)
	return s.db.Exec(sqls, values...)
}
func (s *Session) Model(value interface{}) *Session {
	if s.paresTable == nil || reflect.ValueOf(s.paresTable.Model) != reflect.ValueOf(value) {
		s.paresTable = schema.Parse(value, s.drive)

		if s.tableName != "" {
			s.paresTable.Name = s.tableName
		}
	}
	return s
}
func (s *Session) First(value interface{}) *Session {
	destSlice := reflect.Indirect(reflect.ValueOf(value))
	table := s.Model(value).paresTable
	s.clause.Set(clause.SELECT, table.Name, table.FieldsName)
	s.clause.Set(clause.LIMIT, 1)
	sqls, values := s.clause.Build(clause.SELECT, clause.WHERE, clause.LIMIT)
	row, err := s.db.Query(sqls, values...)
	if err != nil {
		s.err = err
		return s
	}
	for row.Next() {
		var val []interface{}
		for _, name := range table.FieldsName {
			val = append(val, destSlice.FieldByName(name).Addr().Interface())
		}
		if err := row.Scan(val); err != nil {
			s.err = err
			break
		}
	}
	return s
}
func (s *Session) Err() error {
	return s.err
}
func (s *Session) New() *Session {
	session := &Session{
		db:         s.db,
		drive:      s.drive,
		clause:     &clause.Clause{},
		paresTable: nil,
		err:        nil,
	}

	return session
}
func (s *Session) SetTableName(name string) {
	s.tableName = name
}
func (s *Session) Where(dest ...interface{}) {
	s.clause.Set(clause.WHERE, dest...)
	s.clause.Add(clause.WHERE)
}
