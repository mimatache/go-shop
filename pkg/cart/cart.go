package cart

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/logger"

	cartHTTP "github.com/mimatache/go-shop/pkg/cart/http"
	"github.com/mimatache/go-shop/pkg/cart/store"
)

type API struct {
	shoppingCart *cartHTTP.ShoppingCart
}

func NewAPI(logger logger.Logger, inventory cartHTTP.InventoryAPI, users cartHTTP.ClientAPI, payments cartHTTP.PaymentsAPI) (*API, error) {

	cartStore, err := store.New(logger)
	if err != nil {
		return nil, err
	}
	shoppingCart := cartHTTP.New(inventory, users, payments, cartStore)

	return &API{
		shoppingCart: shoppingCart,
	}, nil
}

func (u *API) RegisterToRouter(router *mux.Router, handlers ...func(http.Handler) http.Handler) {
	u.shoppingCart.AddRoutes(router, handlers...)
}
