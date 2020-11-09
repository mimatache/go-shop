package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/http/authorization"
	"github.com/mimatache/go-shop/internal/http/middleware"
	"github.com/mimatache/go-shop/internal/store"
	userStore "github.com/mimatache/go-shop/pkg/users/store"
)

type Authentication interface {
	AddRoutes(router *mux.Router)
}

type userAuthentication struct {
	store userStore.UserStore
}

func New(store userStore.UserStore) Authentication {
	return &userAuthentication{store: store}
}

func (u *userAuthentication) login() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		username, providedPassword, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		user, err := u.store.GetUserByEmail(username)
		// checking error type to evaluate the status code to return
		switch err.(type) {
		case nil:
			// continue with program execution
		case store.NotFound:
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if providedPassword != user.Password {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		tokenString, expirationTime, err := authorization.GenerateJWTToken(user.Name, user.ID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authorization.SetAuthCookie(w, tokenString, expirationTime)

	}
	return http.HandlerFunc(fn)

}

func (u *userAuthentication) logout() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token, err := authorization.GetAuthToken(r)
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		authorization.BlacklistToken(token)
	}

	return http.HandlerFunc(fn)
}

func (u *userAuthentication) AddRoutes(router *mux.Router) {
	router.Handle("/login", u.login()).Methods(http.MethodGet)
	router.Handle("/logout", middleware.JWTAuthorization(u.logout())).Methods(http.MethodGet)
}
