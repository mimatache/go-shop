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

// Validate checks that a user adheres to constraints
func (u Product) Validate() error {
	var errs errors
	if u.ID == 0 {
		errs = append(errs, fmt.Errorf("Product ID cannot be 0"))
	}
	if u.Name == "" {
		errs = append(errs, fmt.Errorf("Name is mandatory"))
	}
	if len(errs) > 0 {
		return errs
	}
	return nil

}
