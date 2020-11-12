package inventory

import (
	"sync"

	"github.com/mimatache/go-shop/pkg/products/store"
)

//go:generate mockgen -source ./inventory.go -destination mocks/inventory.go

// UnderlyingStore represents the interface the invetory store has to implement
type UnderlyingStore interface {
	GetProductByID(ID uint) (*store.Product, error)
	SetProducts(products ...*store.Product) (*store.ProductTransaction, error)
}

// New returns a new instance of inventory
func New(store UnderlyingStore) *Inventory {
	return &Inventory{stock: store}
}

// Inventory represents methods to manage the inventory
type Inventory struct {
	stock UnderlyingStore
	sync.RWMutex
}

// GetProductStock returns the stock of a item given the ID
func (i *Inventory) GetProductStock(id uint) (uint, error) {
	product, err := i.stock.GetProductByID(id)
	if err != nil {
		return 0, err
	}
	return product.GetStock(), nil
}

// HasInStock check if enough stock of a product is still present
func (i *Inventory) HasInStock(id uint, quantity uint) (bool, error) {
	val, err := i.GetProductStock(id)
	if err != nil {
		return false, err
	}
	return val >= quantity, err
}

// GetPrice returns the price of a product
func (i *Inventory) GetPrice(productID uint) (uint, error) {
	product, err := i.stock.GetProductByID(productID)
	if err != nil {
		return 0, err
	}
	return product.GetPrice(), nil
}

// RemoveFromStock removes the requested quantity for each product from stock
// Blocks until the stock is verified as suficient. Blocks writing to store until the condition is met
func (i *Inventory) RemoveFromStock(items map[uint]uint, commitChan <-chan bool, errorChan chan<- error) error {
	i.Lock()
	defer i.Unlock()
	products := []*store.Product{}
	previousProducts := []store.Product{}
	for prodID, desiredQuantity := range items {
		product, err := i.stock.GetProductByID(prodID)
		if err != nil {
			return err
		}
		previousProducts = append(previousProducts, *product)
		err = product.DecreaseStock(desiredQuantity)
		if err != nil {
			return err
		}
		products = append(products, product)
	}
	transaction, err := i.stock.SetProducts(products...)
	if err != nil {
		errorChan <- nil
		return err
	}
	go func() {
		commit := <-commitChan
		if !commit {
			for i := range previousProducts {
				err := transaction.Write(&previousProducts[i])
				if err != nil {
					transaction.Abort()
					errorChan <- err
					return
				}
			}
		}
		transaction.Commit()
		errorChan <- nil
	}()
	return nil
}
