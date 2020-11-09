package store

import (
	"bytes"
	"fmt"
	"regexp"
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

// Email is used to store email addresses and validate them
type Email string

// Validate verifies if the email address is correct
func (e Email) Validate() error {
	ok, err := regexp.Match("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$", []byte(e))
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("invalid email address: %s", e)
	}
	return nil
}

// User models a shop user
type User struct {
	ID       uint   `json:"ID"`
	Name     string `json:"Name"`
	Password string `json:"Password"`
	Email    Email  `json:"Email"`
}

// Validate checks that a user adheres to constraints
func (u User) Validate() error {
	var errs errors
	if u.ID == 0 {
		errs = append(errs, fmt.Errorf("User ID cannot be 0"))
	}
	if u.Name == "" {
		errs = append(errs, fmt.Errorf("Name is mandatory"))
	}
	if u.Password == "" {
		errs = append(errs, fmt.Errorf("Password is mandatory"))
	}
	if err := u.Email.Validate(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errs
	}
	return nil

}
