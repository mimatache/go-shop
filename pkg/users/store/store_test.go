package store_test

import (
	// "fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/mimatache/go-shop/internal/store"
	"github.com/mimatache/go-shop/internal/logger"
	userStore "github.com/mimatache/go-shop/pkg/users/store"
	mock_store "github.com/mimatache/go-shop/pkg/users/store/mocks"
)

const (
	userEmail = "user@email.com"
	password = "1234"
)

var (
	user = &userStore.User{
		Name: "user",
		Email: userEmail,
		Password: password,
		ID: 1,
	}
)

func TestUserStore_GetPasswordByEmail(t *testing.T) {
	g := NewWithT(t)

	log, _, _ := logger.New("test", true)
	

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockUnderlyingStore(ctrl)

	users := userStore.New(log, mockStore)

	mockStore.
		EXPECT().
		Read("user", "email", userEmail).
		Return(user, nil)

	pass, err := users.GetPasswordFor(userEmail)
	
	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(pass).To(Equal(password))
}

func TestUserStore_GiveInvalidEmail(t *testing.T) {
	g := NewWithT(t)

	log, _, _ := logger.New("test", true)
	

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockUnderlyingStore(ctrl)

	users := userStore.New(log, mockStore)

	mockStore.
		EXPECT().
		Read("user", "email", userEmail).
		Return(nil, store.NewNotFoundError("user", "email", userEmail))

	_, err := users.GetPasswordFor(userEmail)
	
	g.Expect(err).Should(HaveOccurred())
}
