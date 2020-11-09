package store_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/mimatache/go-shop/pkg/users/store"
)

func TestEmail_ValidAddresses(t *testing.T) {

	validEmails := []store.Email{
		"someone@email.com",
		"someone_else@email.com",
		"someone.else@email.com",
		"someon3@email.com",
		"someon3@email.com1.com2.com3.com4",
		"abc-d@mail.com",
	}

	for _, email := range validEmails {
		t.Run(string(email), func(t *testing.T) {
			g := NewWithT(t)
			err := email.Validate()
			g.Expect(err).ShouldNot(HaveOccurred(), "email validation reported a good email as bad")
		})
	}

}

func TestEmail_InvalidAddresses(t *testing.T) {

	validEmails := []store.Email{
		"someone@email.com ",
		" someone_else@email.com",
		"someone .else@email.com",
		"abc-@mail.com	",
		"abc..def@mail.com	",
		".abc@mail.com	",
		"abc#def@mail.com",
	}

	for _, email := range validEmails {
		t.Run(string(email), func(t *testing.T) {
			g := NewWithT(t)
			err := email.Validate()
			g.Expect(err).Should(HaveOccurred(), "email validation reported a bad email as good")
		})
	}

}
