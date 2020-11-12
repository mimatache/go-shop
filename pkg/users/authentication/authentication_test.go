package authentication_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/mimatache/go-shop/internal/store"
	"github.com/mimatache/go-shop/pkg/users/authentication"
	mock_authentication "github.com/mimatache/go-shop/pkg/users/authentication/mocks"
)

const (
	goodUser    = "user@email.com"
	invalidUser = "baduser@mail.com"
	goodPasswd  = "testpassword"
	badPassword = "badpassword"
)

func TestUser_ValidCredentials(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	registry := mock_authentication.NewMockUserRegistry(ctrl)

	registry.
		EXPECT().
		GetPasswordFor(goodUser).
		Return(goodPasswd, nil)

	users := authentication.New(registry)

	err := users.IsValid(goodUser, goodPasswd)

	g.Expect(err).ShouldNot(HaveOccurred(), "valid user password combo returned an error")
}

func TestUser_InvalidUsername(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	registry := mock_authentication.NewMockUserRegistry(ctrl)

	registry.
		EXPECT().
		GetPasswordFor(invalidUser).
		Return("", store.NewNotFoundError("users", "name", "name"))

	users := authentication.New(registry)

	err := users.IsValid(invalidUser, goodPasswd)

	g.Expect(err).Should(HaveOccurred(), "invalid user password combo did not return an error")
	g.Expect(err).Should(Equal(authentication.NewInvalidCredentials(invalidUser)))
}

func TestUser_InvalidPassword(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	registry := mock_authentication.NewMockUserRegistry(ctrl)

	registry.
		EXPECT().
		GetPasswordFor(goodUser).
		Return(goodPasswd, nil)

	users := authentication.New(registry)

	err := users.IsValid(goodUser, badPassword)

	g.Expect(err).Should(HaveOccurred(), "invalid user password combo did not return an error")
	g.Expect(err).Should(Equal(authentication.NewInvalidCredentials(goodUser)))
}

func TestInvalidCredentials(t *testing.T) {
	g := NewWithT(t)

	err := fmt.Errorf("some error")

	g.Expect(authentication.IsInvalidCredentialsError(err)).To(BeFalse())

	err = authentication.NewInvalidCredentials("user")
	g.Expect(authentication.IsInvalidCredentialsError(err)).To(BeTrue())
}
