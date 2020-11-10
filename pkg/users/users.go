package users

import (
	"io"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/logger"
	"github.com/mimatache/go-shop/pkg/users/authentication"
	"github.com/mimatache/go-shop/pkg/users/http"
	"github.com/mimatache/go-shop/pkg/users/store"
)

// NewAPI instantiates a new user API and storage
func NewAPI(log logger.Logger, router *mux.Router, db store.UnderlyingStore, seed io.Reader) (*authentication.User, error) {
	users, err := store.New(log, seed, db)
	if err != nil {
		log.Errorf("could not instantiate users DB")
		return nil, err
	}
	authentication := authentication.New(users)
	webAPI := http.New(authentication)
	webAPI.RegisterToRouter(router)
	return authentication, nil
}
