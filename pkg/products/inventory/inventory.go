package inventory

import (
	"fmt"
	"sync"

	db "github.com/mimatache/go-shop/internal/store"

	"github.com/mimatache/go-shop/pkg/products/store"
)

// New returns a new instance of inventory
func New(store store.ProductStore) *Inventory {
	return &Inventory{stock: store}
}

// Inventory represents methods to manage the inventory
type Inventory struct {
	stock store.ProductStore
	sync.RWMutex
}

// GetProductStock returns the stock of a item given the ID
func (a *Inventory) GetProductStock(ID uint) (uint, error) {
	product, err := a.stock.GetProductByID(ID)
	if err != nil {
		return 0, err
	}
	return product.Stock, nil
}

func (a *Inventory) HasInStock(ID uint, quantity uint) (bool, error) {
	val, err := a.GetProductStock(ID)
	if err != nil {
		return false, err
	}
	return val >= quantity, err
}

func (a *Inventory) GetPrice(productID uint) (uint, error) {
	product, err := a.stock.GetProductByID(productID)
	if err != nil {
		return 0, err
	}
	return product.Price, nil
}

// RemoveFromStock removes the requested quantity for each product from stock
// Blocks untill the stock is verified as suficient. Blocks writing to store until the condition is met
func (a *Inventory) RemoveFromStock(items map[uint]uint) (db.ConditionMet, error) {
	a.Lock()
	defer a.Unlock()
	products := []*store.Product{}
	for prodID, desiredQuantity := range items {
		product, err := a.stock.GetProductByID(prodID)
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
