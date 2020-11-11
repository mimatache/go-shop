package store

import (
	"encoding/json"
	"io"

	"github.com/hashicorp/go-memdb"
)

//go:generate mockgen -source ./store.go -destination mocks/store.go

type logger interface {
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
}

var (
	table = &UserTable{name: "user"}
)

// UserTable returns the table for users
type UserTable struct {
	name string
}

// GetName returns the name of the table
func (u *UserTable) GetName() string {
	return u.name
}

// GetTableSchema returns the table schema for users
func (u *UserTable) GetTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: u.name,
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:    "id",
				Unique:  true,
				Indexer: &memdb.UintFieldIndex{Field: "ID"},
			},
			"name": &memdb.IndexSchema{
				Name:    "name",
				Unique:  false,
				Indexer: &memdb.StringFieldIndex{Field: "Name"},
			},
			"password": &memdb.IndexSchema{
				Name:    "password",
				Unique:  false,
				Indexer: &memdb.StringFieldIndex{Field: "Password"},
			},
			"email": &memdb.IndexSchema{
				Name:    "email",
				Unique:  true,
				Indexer: &memdb.StringFieldIndex{Field: "Email"},
			},
		},
	}
}

// UnderlyingStore represents the interface that the DB should implement to be usable
type UnderlyingStore interface {
	Read(table string, key string, value interface{}) (interface{}, error)
	Write(table string, value ...interface{}) error
}

// GetTable returns the user table for the schema
func GetTable() *UserTable {
	return table
}

// LoadSeeds write the seed information to the DB.
func LoadSeeds(seed io.Reader, db UnderlyingStore) error {
	var users []*User

	err := json.NewDecoder(seed).Decode(&users)
	if err != nil {
		return err
	}
	for _, user := range users {
		err = db.Write(table.GetName(), user)
		if err != nil {
			return err
		}
	}
	return nil
}

// New creates a new user DB instance
func New(log logger, db UnderlyingStore) UserStore {
	return &userLogger{
		log:   log,
		store: &userStore{db: db},
	}
}

// UserStore models the user DB
type UserStore interface {
	// GetPasswordFor returns the password for the given user
	GetPasswordFor(name string) (string, error)
}

type userStore struct {
	db UnderlyingStore
}

func (u *userStore) GetPasswordFor(email string) (string, error) {
	user, err := checkAndReturn(u.db.Read(table.GetName(), "email", email))
	if err != nil {
		return "", err
	}
	return user.Password, nil
}

type userLogger struct {
	log   logger
	store UserStore
}

func (u *userLogger) GetPasswordFor(name string) (string, error) {

	var err error
	defer func() {
		if err != nil {
			u.log.Debugf("error occured when retrieving password for user %s", name)
			u.log.Debugf("%v", err)
			return
		}
		u.log.Debugf("Retrieved password for user %s", name)
	}()
	user, err := u.store.GetPasswordFor(name)
	return user, err
}

// checkAndReturn reads the output from the DB and returns a User instance if no error occured.
// This will panic if the DB does not return and error but the output is not an User. 
// Intentinally left to do this as if this happens it means we have an incosistency in the DB that should be resolve immediately
// and silent handling might mask this issue
func checkAndReturn(raw interface{}, err error) (*User, error) {
	if err != nil {
		return nil, err
	}
	return raw.(*User), nil
}
