package users

import (
	"io"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/logger"
	"github.com/mimatache/go-shop/pkg/users/http"
	"github.com/mimatache/go-shop/pkg/users/store"
)

func NewAPI(log logger.Logger, seed io.Reader) (*API, error) {
	users, err := store.New(log, seed)
	if err != nil {
		log.Errorf("could not instantiate users DB")
		return nil, err
	}
	user := http.New(users)
	return &API{
		http:  user,
		store: users,
	}, nil
}

type API struct {
	http  http.Authentication
	store store.UserStore
}

func (u *API) RegisterToRouter(router *mux.Router) {
	u.http.AddRoutes(router)
}

func (u *API) GetEmailForUser(ID uint) (string, error) {
	user, err := u.store.GetUserByID(ID)
	if err != nil {
		return "", err
	}
	return string(user.Email), nil
}
