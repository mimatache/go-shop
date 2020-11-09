package store

import (
	"fmt"
	"github.com/hashicorp/go-memdb"
)

type ConditionMet func(write bool)


type NotFound struct {
	Msg string
}

func (u NotFound) Error() string {
	return u.Msg
}

type Store interface {
	Write(table Table, obj ...interface{}) error
	WriteAfterExternalCondition(table Table, objs ...interface{}) (ConditionMet, error)
	Read(table Table, key string, value interface{}) (interface{}, error)
	ReadAll(table Table, key string, value interface{}) ([]interface{}, error)
	Remove(table Table, key string, value interface{}) error
}

func New(schema Schema) (Store, error) {
	db, err := schema.initDB()
	if err != nil {
		return nil, err
	}
	return &store{db: db}, nil
}

type store struct {
	db *memdb.MemDB
}

func (s *store) Write(table Table, objs ...interface{}) error {
	finalizer, err := s.WriteAfterExternalCondition(table, objs...)
	if err != nil {
		return err
	}
	finalizer(true)
	return nil
}

func (s *store) WriteAfterExternalCondition(table Table, objs ...interface{}) (ConditionMet, error) {
	txn := s.db.Txn(true)
	shouldCommit := s.shouldCommit(txn, s.db.Snapshot())
	for _, obj := range objs {
		if err := txn.Insert(table.GetName(), obj); err != nil {
			txn.Abort()
			return nil, err
		}
	}
	return shouldCommit, nil
}

func (s *store) Read(table Table, key string, value interface{}) (interface{}, error) {
	txn := s.db.Txn(false)
	raw, err := txn.First(table.GetName(), key, value)

	if raw == nil {
		return nil, NotFound{fmt.Sprintf("could not find %s with %s equal to %v", table.GetName(), key, value)}
	}
	return raw, err
}

func (s *store) ReadAll(table Table, key string, value interface{}) ([]interface{}, error) {
	txn := s.db.Txn(false)
	it, err := txn.Get(table.GetName(), key, value)
	if err != nil {
		return nil, err
	}
	items := []interface{}{}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		items = append(items, obj)
	}
	return items, nil
}

func (s *store) Remove(table Table, key string, value interface{}) error {
	txn := s.db.Txn(true)
	_, err := txn.DeleteAll(table.GetName(), key, value)
	if err != nil {
		txn.Abort()
		return err
	}
	txn.Commit()
	return nil
}


func (s *store) shouldCommit(txn *memdb.Txn, db *memdb.MemDB) ConditionMet {
	return func(write bool) {
		if write {
			txn.Commit()
		} else {
			txn.Abort()
			s.db = db
		}
	}
}
