package cart

import (
	netHTTP "net/http"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/logger"

	"github.com/mimatache/go-shop/pkg/cart/http"
	"github.com/mimatache/go-shop/pkg/cart/store"
)

// NewAPI instantiates a new cart API
func NewAPI(
	logger logger.Logger, 
	inventory http.InventoryAPI,
	users http.ClientAPI, 
	payments http.PaymentsAPI, 
	db store.UnderlyingStore,
	router *mux.Router, 
	handlers ...func(netHTTP.Handler) netHTTP.Handler,
) error {
	cartStore, err := store.New(logger, db)
	if err != nil {
		return err
	}
	shoppingCart := http.New(inventory, users, payments, cartStore)
	shoppingCart.AddRoutes(router, handlers...)

	return nil
}

