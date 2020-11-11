package store_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/mimatache/go-shop/internal/logger"
	"github.com/mimatache/go-shop/internal/store"
	productStore "github.com/mimatache/go-shop/pkg/products/store"
	mock_store "github.com/mimatache/go-shop/pkg/products/store/mocks"
)

const (
	productID uint = 1
)

var (
	item = &productStore.Product{
		ID: productID,
		Name: "awesome product",
		Price: 10,
		Stock: 2,
	}

	table = productStore.GetTable().GetName()
)

func TestProductStore_GetProductByID(t *testing.T) {
	g := NewWithT(t)

	log, _, _ := logger.New("test", true)
	

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockUnderlyingStore(ctrl)

	products := productStore.New(log, mockStore)

	mockStore.
		EXPECT().
		Read(table, "id", productID).
		Return(item, nil)

	stockItem, err := products.GetProductByID(productID)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(stockItem).To(Equal(item))
}

func TestProductStore_GetProductByID_NotFound(t *testing.T) {
	g := NewWithT(t)

	log, _, _ := logger.New("test", true)
	

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockUnderlyingStore(ctrl)

	products := productStore.New(log, mockStore)

	mockStore.
		EXPECT().
		Read(table, "id", productID).
		Return(nil, store.NewNotFoundError(table, "id", productID))

	_, err := products.GetProductByID(productID)

	g.Expect(err).Should(HaveOccurred())
}

func TestProductStore_SetProducts(t *testing.T) {
	g := NewWithT(t)

	log, _, _ := logger.New("test", true)
	

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockUnderlyingStore(ctrl)

	products := productStore.New(log, mockStore)

	mockStore.
		EXPECT().
		WriteAfterExternalCondition(table, item).
		Return(func(bool) {} ,nil)

	_, err := products.SetProducts(item)

	g.Expect(err).ShouldNot(HaveOccurred())
}

func TestProductStore_SetProducts_Error(t *testing.T) {
	g := NewWithT(t)

	log, _, _ := logger.New("test", true)
	

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_store.NewMockUnderlyingStore(ctrl)

	products := productStore.New(log, mockStore)

	mockStore.
		EXPECT().
		WriteAfterExternalCondition(table, item).
		Return(nil ,fmt.Errorf("random error"))

	_, err := products.SetProducts(item)

	g.Expect(err).Should(HaveOccurred())
}