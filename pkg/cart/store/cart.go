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

type CartItem struct {
	ID       uint          `json:"id"`
	Products map[uint]uint `json:"products"`
}

func (c CartItem) Validate() error {
	var errs errors
	if c.ID == 0 {
		errs = append(errs, fmt.Errorf("User ID cannot be 0"))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func NewCartItem(user uint, product uint, quantity uint) (*CartItem, error) {
	c := CartItem{
		ID:       user,
		Products: map[uint]uint{product: quantity},
	}
	err := c.Validate()
	if err != nil {
		return nil, err
	}
	return &c, nil
}
