package products

import (
	"fmt"
	"io"
	"sync"

	"github.com/mimatache/go-shop/internal/logger"
	"github.com/mimatache/go-shop/internal/store"
	productsStore "github.com/mimatache/go-shop/pkg/products/store"
)

func NewAPI(log logger.Logger, seed io.Reader) (*API, error) {
	stock, err := productsStore.New(log, seed)
	if err != nil {
		log.Errorf("could not instantiate products DB")
		return nil, err
	}
	return &API{stock: stock}, nil
}

type API struct {
	stock productsStore.ProductStore
	sync.RWMutex
}

// GetProductStock returns the stock of a item given the ID
func (a *API) GetProductStock(ID uint) (uint, error) {
	product, err := a.stock.GetProductById(ID)
	if err != nil {
		return 0, err
	}
	return product.Stock, nil
}

func (a *API) HasInStock(ID uint, quantity uint) (bool, error) {
	val, err := a.GetProductStock(ID)
	if err != nil {
		return false, err
	}
	return val >= quantity, err
}

func (a *API) GetPrice(productID uint) (uint, error) {
	product, err := a.stock.GetProductById(productID)
	if err != nil {
		return 0, err
	}
	return product.Price, nil
}

// RemoveFromStock removes the requested quantity for each product from stock
// Blocks untill the stock is verified as suficient. Blocks writing to store until the condition is met
func (a *API) RemoveFromStock(items map[uint]uint) (store.ConditionMet, error) {
	a.Lock()
	defer a.Unlock()
	products := []*productsStore.Product{}
	for prodID, desiredQuantity := range items {
		product, err := a.stock.GetProductById(prodID)
		if err != nil {
			return nil, err
		}
		if product.Stock < desiredQuantity {
			return nil, fmt.Errorf("insuficient stock of %s", product.Name)
		}
		product.Stock -= desiredQuantity
		products = append(products, product)
	}
	isConditionMet, err := a.stock.SetProducts(products...)
	if err != nil {
		return nil, fmt.Errorf("could not remove products from stock")
	}
	return isConditionMet, nil
}
