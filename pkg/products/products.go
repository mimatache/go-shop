package products

import (
	"io"

	"github.com/mimatache/go-shop/internal/logger"
	"github.com/mimatache/go-shop/pkg/products/inventory"
	"github.com/mimatache/go-shop/pkg/products/store"
)

func NewAPI(log logger.Logger, db store.UnderlyingStore, seed io.Reader) (*inventory.Inventory, error) {
	stock, err := store.New(log, db, seed)
	if err != nil {
		log.Errorf("could not instantiate products DB")
		return nil, err
	}

	return inventory.New(stock), nil
}
