package store

import (
	"bytes"
	"fmt"
)

type errors []error

func (e errors) Error() string {
	b := bytes.NewBufferString("")
	for _, err := range e {
		_, _ = fmt.Fprintf(b, "\t%s\n", err)
	}
	return b.String()
}

// Validatable is an item which has to adhere to certain convetions and knows how to check that it is correct
type Validatable interface {
	Validate() error
}

// Product models a shop product
type Product struct {
	ID    uint   `json:"ID"`
	Name  string `json:"Name"`
	Price uint   `json:"Price"`
	Stock uint   `json:"Stock"`
}

// GetPrice returns the amount of this product left in stock
func (p *Product) GetPrice() uint {
	return p.Price
}

// GetStock returns the amount of this product left in stock
func (p *Product) GetStock() uint {
	return p.Stock
}

// HasStock returns true is the stock is higher or equal to the requeted ammount
func (p *Product) HasStock(quantity uint) bool {
	return p.Stock >= quantity
}

// IncreaseStock adds the given quantitty to the stock o
func (p *Product) IncreaseStock(quantity uint) {
	p.Stock += quantity
}

// DecreaseStock decreases the quantity of an item if sufficient
func (p *Product) DecreaseStock(quantitty uint) error {
	if p.HasStock(quantitty) {
		p.Stock -= quantitty
		return nil
	}
	return fmt.Errorf("insuficient stock of %s", p.Name)
}

// Validate checks that a user adheres to constraints
func (p *Product) Validate() error {
	var errs errors
	if p.ID == 0 {
		errs = append(errs, fmt.Errorf("Product ID cannot be 0"))
	}
	if p.Name == "" {
		errs = append(errs, fmt.Errorf("Name is mandatory"))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil

}
