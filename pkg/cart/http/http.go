package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/http/authorization"
	"github.com/mimatache/go-shop/internal/http/helpers"
	"github.com/mimatache/go-shop/pkg/cart/cart"
)

func New(cart *cart.Cart) *ShoppingCart {
	return &ShoppingCart{
		cart: cart,
	}
}

type ShoppingCart struct {
	cart *cart.Cart
}

func (s *ShoppingCart) addProductToCart(w http.ResponseWriter, r *http.Request) {
	userID, err := authorization.GetUserIDFromRequest(r)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var prod cart.Product
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&prod)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentContents, err := s.cart.AddProductToCart(userID, prod)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusBadRequest)
		return
	}

	helpers.FormatResponse(w, currentContents)
}

func (s *ShoppingCart) checkout(w http.ResponseWriter, r *http.Request) {
	userID, err := authorization.GetUserIDFromRequest(r)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusUnauthorized)
		return
	}
	contents, err := s.cart.Checkout(userID)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	helpers.FormatResponse(w, contents)
}

// AddRoutes registers the API routes to a router
func (s *ShoppingCart) AddRoutes(router *mux.Router, handlers ...func(http.Handler) http.Handler) {
	cartRouter := router.PathPrefix("/cart").Subrouter()
	cartRouter.HandleFunc("/add", s.addProductToCart).Methods(http.MethodPost)
	cartRouter.HandleFunc("/checkout", s.checkout).Methods(http.MethodPost)
	for _, v := range handlers {
		cartRouter.Use(v)
	}
}
