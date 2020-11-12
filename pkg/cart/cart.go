package cart

import (
	netHTTP "net/http"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/logger"

	"github.com/mimatache/go-shop/pkg/cart/cart"
	"github.com/mimatache/go-shop/pkg/cart/http"
	"github.com/mimatache/go-shop/pkg/cart/store"
)

// NewAPI instantiates a new cart API
func NewAPI(
	logger logger.Logger,
	inventory cart.InventoryAPI,
	payments cart.PaymentsAPI,
	db store.UnderlyingStore,
	router *mux.Router,
	handlers ...func(netHTTP.Handler) netHTTP.Handler,
) {
	cartStore := store.New(logger, db)
	cart := cart.New(inventory, payments, cartStore)
	shoppingCart := http.New(cart)
	shoppingCart.AddRoutes(router, handlers...)
}
