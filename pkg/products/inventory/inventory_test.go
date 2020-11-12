package inventory_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/mimatache/go-shop/pkg/products/inventory"
	mock_inventory "github.com/mimatache/go-shop/pkg/products/inventory/mocks"
	"github.com/mimatache/go-shop/pkg/products/store"
)

const (
	itemID uint = 1
	stock  uint = 3
	price  uint = 100
)

var (
	product = store.Product{
		ID:    itemID,
		Name:  "Product 1",
		Price: price,
		Stock: stock,
	}
)

func TestInventory_GetProductStock(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInventory := mock_inventory.NewMockUnderlyingStore(ctrl)

	productInventory := inventory.New(mockInventory)

	mockInventory.
		EXPECT().
		GetProductByID(itemID).
		Return(&product, nil)

	itemStock, err := productInventory.GetProductStock(itemID)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(itemStock).To(Equal(uint(stock)))

}

func TestInventory_GetProductStock_Error(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInventory := mock_inventory.NewMockUnderlyingStore(ctrl)

	productInventory := inventory.New(mockInventory)

	mockInventory.
		EXPECT().
		GetProductByID(itemID).
		Return(nil, fmt.Errorf("an error"))

	_, err := productInventory.GetProductStock(itemID)

	g.Expect(err).Should(HaveOccurred())
}

func TestInventory_HasInStock(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInventory := mock_inventory.NewMockUnderlyingStore(ctrl)

	productInventory := inventory.New(mockInventory)

	mockInventory.
		EXPECT().
		GetProductByID(itemID).
		Return(&product, nil)

	has, err := productInventory.HasInStock(itemID, stock)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(has).To(BeTrue())

}

func TestInventory_HasInStock_False(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInventory := mock_inventory.NewMockUnderlyingStore(ctrl)

	productInventory := inventory.New(mockInventory)

	mockInventory.
		EXPECT().
		GetProductByID(itemID).
		Return(&product, nil)

	has, err := productInventory.HasInStock(itemID, stock+1)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(has).To(BeFalse())
}

func TestInventory_HasInStock_Error(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInventory := mock_inventory.NewMockUnderlyingStore(ctrl)

	productInventory := inventory.New(mockInventory)

	mockInventory.
		EXPECT().
		GetProductByID(itemID).
		Return(nil, fmt.Errorf("an error"))

	_, err := productInventory.HasInStock(itemID, stock+1)

	g.Expect(err).Should(HaveOccurred())
}

func TestInventory_GetPrice(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInventory := mock_inventory.NewMockUnderlyingStore(ctrl)

	productInventory := inventory.New(mockInventory)

	mockInventory.
		EXPECT().
		GetProductByID(itemID).
		Return(&product, nil)

	productPrice, err := productInventory.GetPrice(itemID)

	g.Expect(err).ShouldNot(HaveOccurred())
	g.Expect(productPrice).To(Equal(price))

}

func TestInventory_GetPrice_Error(t *testing.T) {
	g := NewWithT(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInventory := mock_inventory.NewMockUnderlyingStore(ctrl)

	productInventory := inventory.New(mockInventory)

	mockInventory.
		EXPECT().
		GetProductByID(itemID).
		Return(nil, fmt.Errorf("an error"))

	_, err := productInventory.GetPrice(itemID)

	g.Expect(err).Should(HaveOccurred())
}
