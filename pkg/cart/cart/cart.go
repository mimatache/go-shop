package cart

import (
	"bytes"
	"fmt"

	"github.com/mimatache/go-shop/internal/store"
	shoppingCart "github.com/mimatache/go-shop/pkg/cart/store"
)

type errors []error

func (e errors) Error() string {
	b := bytes.NewBufferString("")
	for _, err := range e {
		_, _ = fmt.Fprintf(b, "\t%s\n", err)
	}
	return b.String()
}

// InventoryAPI represents the methods that need to be implemented by the inventory API
type InventoryAPI interface {
	// Check if product is in sufficient stock
	HasInStock(productID uint, quantity uint) (bool, error)
	// GetPrice returns the price of an item
	GetPrice(productID uint) (uint, error)
	// Renove from stock removes items from stock, but restores the previous values of False is sent over the commitChan
	RemoveFromStock(items map[uint]uint, commitChan <-chan bool, errorChan chan<- error) error
}

// PaymentsAPI represents the methods that need to be implemented by the payments API
type PaymentsAPI interface {
	// MakePayment tries to call the payment API
	MakePayment(client string, money uint) error
}

// Product represents a product added to the cart
type Product struct {
	ID       uint `json:"id"`
	Quantity uint `json:"quantity"`
}

// Contents represents the contents of the cart
type Contents struct {
	Products []*Product `json:"products"`
}

// New starts a new cart
func New(inventory InventoryAPI, payments PaymentsAPI, cartContents shoppingCart.CartStore) *Cart {
	return &Cart{
		inventory:    inventory,
		payments:     payments,
		cartContents: cartContents,
	}
}

// Cart represents the actions that you can perform on a cart
type Cart struct {
	inventory    InventoryAPI
	payments     PaymentsAPI
	cartContents shoppingCart.CartStore
}

// Checkout attempts to perform checkout of the current cart contents
func (c *Cart) Checkout(userID string) (*Contents, error) {
	cartContents, err := c.getContents(userID)
	if err != nil {
		return nil, err
	}

	var cost uint
	items := map[uint]uint{}
	for _, item := range cartContents.Products {
		price, err := c.inventory.GetPrice(item.ID)
		if err != nil {
			return nil, err
		}
		cost += price * item.Quantity
		items[item.ID] = item.Quantity
	}

	errChan := make(chan error)
	commitChan := make(chan bool)

	defer close(errChan)
	defer close(commitChan)

	err = c.inventory.RemoveFromStock(items, commitChan, errChan)
	if err != nil {
		return nil, err
	}
	errs := errors{}
	err = c.payments.MakePayment(userID, cost)
	if err != nil {
		commitChan <- false
		errs = append(errs, err)
	} else {
		commitChan <- true
	}
	err = <-errChan
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	err = c.cartContents.ClearCartFor(userID)
	if err != nil {
		return nil, err
	}
	return cartContents, err
}

// AddProductToCart adds a bew product to the cart (or updates the existing item quantity if some already present)
func (c *Cart) AddProductToCart(userID string, prod Product) (*Contents, error) {

	currentProducts, err := c.cartContents.GetProductsForUser(userID)
	if err != nil && !store.IsNotFoundError(err) {
		return nil, err
	}
	quantity := currentProducts[prod.ID]

	hasStock, err := c.inventory.HasInStock(prod.ID, quantity+prod.Quantity)
	if err != nil {
		return nil, err
	}
	if !hasStock {
		return nil, fmt.Errorf("Insuficient stock")
	}

	_, err = c.cartContents.AddProduct(userID, prod.ID, prod.Quantity)
	if err != nil {
		return nil, err
	}

	currentContents, err := c.getContents(userID)
	if err != nil {
		return nil, err
	}

	return currentContents, nil
}

func (c *Cart) getContents(userID string) (*Contents, error) {
	currentProducts, err := c.cartContents.GetProductsForUser(userID)
	if err != nil {
		return nil, err
	}
	currentContents := &Contents{Products: []*Product{}}
	for k, v := range currentProducts {
		currentContents.Products = append(currentContents.Products, &Product{ID: k, Quantity: v})
	}
	return currentContents, nil
}
