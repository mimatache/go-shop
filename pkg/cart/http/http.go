package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/http/authorization"
	"github.com/mimatache/go-shop/internal/http/helpers"
	"github.com/mimatache/go-shop/internal/store"
	shoppingCart "github.com/mimatache/go-shop/pkg/cart/store"
)

type InventoryAPI interface {
	HasInStock(productID uint, quantity uint) (bool, error)
	GetPrice(productID uint) (uint, error)
	RemoveFromStock(items map[uint]uint) (store.ConditionMet, error)
}

type PaymentsAPI interface {
	MakePayment(client string, money uint) error
}

type ClientAPI interface {
	GetEmailForUser(ID uint) (string, error)
}

// Product represents a product added to the cart
type Product struct {
	ID       uint `json:"id"`
	Quantity uint `json:"quantity"`
}

type Contents struct {
	Products []*Product `json:"products"`
}

func New(inventory InventoryAPI, clientAPI ClientAPI, payments PaymentsAPI, cartContents shoppingCart.CartStore) *ShoppingCart {
	return &ShoppingCart{
		inventory:    inventory,
		clientAPI:    clientAPI,
		cartContents: cartContents,
		payments:     payments,
	}
}

type ShoppingCart struct {
	inventory    InventoryAPI
	clientAPI    ClientAPI
	payments     PaymentsAPI
	cartContents shoppingCart.CartStore
}

func (s *ShoppingCart) addProductToCart(w http.ResponseWriter, r *http.Request) {
	userID, err := authorization.GetUserIdFromRequest(r)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var prod Product

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&prod)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusBadRequest)
		return
	}

	isMet, quantity, err := s.cartContents.AddProductWhenConditionIsMet(userID, prod.ID, prod.Quantity)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusBadRequest)
		return
	}
	hasStock, err := s.inventory.HasInStock(prod.ID, quantity)
	isMet(hasStock)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !hasStock {
		helpers.FormatError(w, fmt.Sprintf("insuficient stock of product %d", prod.ID), http.StatusBadRequest)
		return
	}

	currentContents, err := s.getContents(userID)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(currentContents)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func (s *ShoppingCart) checkout(w http.ResponseWriter, r *http.Request) {
	userID, err := authorization.GetUserIdFromRequest(r)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusUnauthorized)
		return
	}
	cartContents, err := s.getContents(userID)
	switch err.(type){
	case nil:
		// nothing to do
	case store.NotFound:
		helpers.FormatError(w, err.Error(), http.StatusNotFound)
		return
	default:
		helpers.FormatError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	var cost uint
	items := map[uint]uint{}
	for _, item := range cartContents.Products {
		price, err := s.inventory.GetPrice(item.ID)
		if err != nil {
			helpers.FormatError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cost += price * item.Quantity
		items[item.ID] = item.Quantity
	}

	isMet, err := s.inventory.RemoveFromStock(items)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusBadRequest)
		return
	}
	userEmail, err := s.clientAPI.GetEmailForUser(userID)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.payments.MakePayment(userEmail, cost)
	if err != nil {
		isMet(false)
		helpers.FormatError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	isMet(true)
	
	err = s.cartContents.ClearCartFor(userID)
	if err != nil {
		helpers.FormatError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *ShoppingCart) AddRoutes(router *mux.Router, handlers ...func(http.Handler) http.Handler) {
	cartRouter := router.PathPrefix("/cart").Subrouter()
	cartRouter.HandleFunc("/add", s.addProductToCart).Methods(http.MethodPost)
	cartRouter.HandleFunc("/checkout", s.checkout).Methods(http.MethodPost)
	for _, v := range handlers {
		cartRouter.Use(v)
	}
}

func (s *ShoppingCart) getContents(userID uint) (*Contents, error) {
	currentProducts, err := s.cartContents.GetProductsForUser(userID)
	if err != nil {
		return nil, err
	}
	currentContents := &Contents{Products: []*Product{}}
	for k, v := range currentProducts {
		currentContents.Products = append(currentContents.Products, &Product{ID: k, Quantity: v})
	}
	return currentContents, nil
}
