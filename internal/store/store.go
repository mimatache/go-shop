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

// NewNotFoundError return a not found error for the given input
func NewNotFoundError(tableName string, key string, value interface{}) error {
	return NotFound{fmt.Sprintf("could not find %s with %s equal to %v", tableName, key, value)}
}

// IsNotFoundError checks if an error is of type notFound
func IsNotFoundError(err error) bool {
	switch err.(type) {
	case NotFound:
		return true
	default:
		return false
	}
}

// New createa a new Store with the given schema
func New(schema Schema) (*Store, error) {
	db, err := schema.initDB()
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Store is used to connect to the database
type Store struct {
	db *memdb.MemDB
}

// Write inserts a row into the table
func (s *Store) Write(table string, objs ...interface{}) error {
	finalizer, err := s.WriteAfterExternalCondition(table, objs...)
	if err != nil {
		return err
	}
	finalizer(true)
	return nil
}

// WriteAfterExternalCondition permanently adds items in the DB only after an external condition is met
func (s *Store) WriteAfterExternalCondition(table string, objs ...interface{}) (ConditionMet, error) {
	txn := s.db.Txn(true)
	shouldCommit := s.shouldCommit(txn, s.db.Snapshot())
	for _, obj := range objs {
		if err := txn.Insert(table, obj); err != nil {
			txn.Abort()
			return nil, err
		}
	}
	return shouldCommit, nil
}

// Read returns a row from a DB table
func (s *Store) Read(table string, key string, value interface{}) (interface{}, error) {
	txn := s.db.Txn(false)
	raw, err := txn.First(table, key, value)

	if raw == nil {
		return nil, NewNotFoundError(table, key, value)
	}
	return raw, err
}

// Remove removes a row from the DB table
func (s *Store) Remove(table string, key string, value interface{}) error {
	txn := s.db.Txn(true)
	_, err := txn.DeleteAll(table, key, value)
	if err != nil {
		txn.Abort()
		return err
	}
	txn.Commit()
	return nil
}

func (s *Store) shouldCommit(txn *memdb.Txn, db *memdb.MemDB) ConditionMet {
	return func(write bool) {
		if write {
			txn.Commit()
		} else {
			txn.Abort()
			s.db = db
		}
	}
}
