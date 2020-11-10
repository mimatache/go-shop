package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/http/authorization"
	"github.com/mimatache/go-shop/internal/http/helpers"
	"github.com/mimatache/go-shop/internal/http/middleware"
	"github.com/mimatache/go-shop/pkg/users/authentication"
)

// AuthenticationAPI is used to authenticate users
type AuthenticationAPI struct {
	users userAuthentication
}

type userAuthentication interface {
	IsValid(email, password string) error
}

// New creates a new AuthenticationApi
func New(users userAuthentication) *AuthenticationAPI {
	return &AuthenticationAPI{users: users}
}

func (u *AuthenticationAPI) login() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		username, providedPassword, ok := r.BasicAuth()
		if !ok {
			helpers.FormatError(w, authentication.NewInvalidCredentials(username).Error(), http.StatusUnauthorized)
			return
		}
		err := u.users.IsValid(username, providedPassword)
		if err != nil {
			if authentication.IsInvalidCredentialsError(err) {
				helpers.FormatError(w, err.Error(), http.StatusUnauthorized)
				return
			}
			helpers.FormatError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tokenString, expirationTime, err := authorization.GenerateJWTToken(username)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authorization.SetAuthCookie(w, tokenString, expirationTime)

	}
	return http.HandlerFunc(fn)

}

func (u *AuthenticationAPI) logout() http.Handler {
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

// RegisterToRouter adds the API routes to the given router
func (u *AuthenticationAPI) RegisterToRouter(router *mux.Router) {
	router.Handle("/login", u.login()).Methods(http.MethodGet)
	router.Handle("/logout", middleware.JWTAuthorization(u.logout())).Methods(http.MethodGet)
}
