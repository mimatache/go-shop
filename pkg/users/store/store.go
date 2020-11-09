package store

import (
	"encoding/json"
	"io"

	"github.com/hashicorp/go-memdb"

	"github.com/mimatache/go-shop/internal/store"
)

type logger interface {
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
}

var (
	table = &userTable{name: "user"}
)

type userTable struct {
	name string
}

func (u *userTable) GetName() string {
	return u.name
}

func (u *userTable) GetTableSchema() *memdb.TableSchema {
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

func New(log logger, seed io.Reader) (UserStore, error) {
	schema := store.NewSchema()

	var users []*User

	err := json.NewDecoder(seed).Decode(&users)
	if err != nil {
		return nil, err
	}

	schema.AddToSchema(table)

	db, err := store.New(schema)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		err = db.Write(table, user)
		if err != nil {
			return nil, err
		}
	}

	return &userLogger{
		log:   log,
		store: &userStore{db: db},
	}, nil
}

type UserStore interface {
	GetUserByName(name string) (*User, error)
	GetUserByID(ID uint) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetPasswordFor(name string) (string, error)
}

type userStore struct {
	db store.Store
}

func checkAndReturn(raw interface{}, err error) (*User, error) {
	if err != nil {
		return nil, err
	}
	return raw.(*User), nil
}

func (u *userStore) GetUserByName(name string) (*User, error) {
	return checkAndReturn(u.db.Read(table, "name", name))
}

func (u *userStore) GetUserByEmail(email string) (*User, error) {
	return checkAndReturn(u.db.Read(table, "email", email))
}

func (u *userStore) GetPasswordFor(name string) (string, error) {
	user, err := u.GetUserByEmail(name)
	if err != nil {
		return "", err
	}
	return user.Password, nil
}

func (u *userStore) GetUserByID(ID uint) (*User, error) {
	return checkAndReturn(u.db.Read(table, "id", ID))
}

type userLogger struct {
	log   logger
	store UserStore
}

func (u *userLogger) GetUserByName(name string) (*User, error) {
	var err error
	defer func() {
		if err != nil {
			u.log.Debugf("error occured when retrieving user %s", name)
			return
		}
		u.log.Debugf("Retrieved user %s", name)
	}()
	user, err := u.store.GetUserByName(name)
	return user, err
}

func (u *userLogger) GetUserByID(ID uint) (*User, error) {
	var err error
	defer func() {
		if err != nil {
			u.log.Debugf("error occured when retrieving user %d", ID)
			return
		}
		u.log.Debugf("Retrieved user %d", ID)
	}()
	user, err := u.store.GetUserByID(ID)
	return user, err
}

func (u *userLogger) GetUserByEmail(name string) (*User, error) {
	var err error
	defer func() {
		if err != nil {
			u.log.Debugf("error occured when retrieving user %s", name)
			return
		}
		u.log.Debugf("Retrieved user %s", name)
	}()
	user, err := u.store.GetUserByEmail(name)
	return user, err
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
