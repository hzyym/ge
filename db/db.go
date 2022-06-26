package db

import (
	"gz/db/inter"
	"gz/db/session"
)

type Db struct {
	session *session.Session

	clone int
}

func Open(drive inter.Drive) (*Db, error) {
	sessions := new(session.Session)
	err := sessions.Open(drive)
	if err != nil {
		return nil, err
	}
	return &Db{
		session: sessions,
	}, nil
}
func (d *Db) Where(dest ...interface{}) *Db {
	tx := d.getDB()
	tx.session.Where(dest...)
	return tx
}
func (d *Db) First(dest interface{}) {
	tx := d.getDB()
	tx.session.First(dest)
}
func (d *Db) Table(name string) *Db {
	tx := d.getDB()
	tx.session.SetTableName(name)
	return tx
}
func (d *Db) getDB() *Db {
	if d.clone == 0 {
		db := &Db{
			session: d.session.New(),
			clone:   1,
		}
		return db
	}
	return d
}
