package products

import (
	"github.com/mimatache/go-shop/internal/logger"
	"github.com/mimatache/go-shop/pkg/products/inventory"
	"github.com/mimatache/go-shop/pkg/products/store"
)

func NewAPI(log logger.Logger, db store.UnderlyingStore) *inventory.Inventory {
	stock := store.New(log, db)
	return inventory.New(stock)
}
