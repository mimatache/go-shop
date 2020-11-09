package client

import (
	"bytes"
	"fmt"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
)

type ShopClient struct{
	Name string
	password string
	client *http.Client
	shopAPI string
}

type Product struct {
	ID uint `json:"id"`
	Quantity uint `json:"quantity"`
}

type ErrorBody struct {
	Error string `json:"error"`
	Code int `json:"code"`
}

type Client interface{
	Login() error
	AddToCart(prodID, quantity uint) error
	Checkout() error
}

type Actor interface{
	PerformActionLoop(prodID, quantity uint) error
}

func New(username, password, shopAPI string) (*ShopClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &ShopClient{
		Name: username,
		password: password,
		shopAPI: shopAPI,
		client: &http.Client{Jar: jar},
	}, nil
}

func (s *ShopClient) Login() error {
	r, err := http.NewRequest(http.MethodGet, s.makeURL("/login"), nil)
	if err != nil {
		return err
	}
	s.basicAuth(r)
	resp, err := s.client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *ShopClient) AddToCart(prodId, qunatity uint) error {
	product := &Product{
		ID: prodId,
		Quantity: qunatity,
	}
	reqBody, err := json.Marshal(product)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, s.makeURL("cart/add"), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
    req.Header.Set("Content-Type", "application/json")
	response, err := s.client.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return handleBadStatusCodes(response)
	}
	response.Body.Close()		
	return nil
}

func (s *ShopClient) Checkout() error {
	req, err := http.NewRequest(http.MethodPost, s.makeURL("cart/checkout"), bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
    req.Header.Set("Content-Type", "application/json")
	response, err := s.client.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return handleBadStatusCodes(response)
	}
	response.Body.Close()		
	return nil
}

func (s *ShopClient) PerformActionLoop(itemsToBuy map[uint]uint) error {
	err := s.Login()
	if err != nil{
		return err
	}
	for prodID, qunatity := range itemsToBuy{
		var i uint
		for i = 0; i < qunatity; i++ {
			err = s.AddToCart(prodID, 1)
			if err != nil{
				return fmt.Errorf("error adding to cart %v", err)
		}
		}
	}
	err = s.Checkout()
	if err != nil{
		return  fmt.Errorf("error adding doing checkout %v", err)
	}
	return nil
}

func (s *ShopClient) makeURL(path string) string{
	return fmt.Sprintf("%s/%s", s.shopAPI, path)
}

func (s *ShopClient) basicAuth(r *http.Request){
	r.Header.Add("Authorization","Basic " + base64.StdEncoding.EncodeToString([]byte(s.Name + ":" + s.password)))
}

func handleBadStatusCodes(r *http.Response) error{
	errorResp := &ErrorBody{}
	if err := json.NewDecoder(r.Body).Decode(errorResp); err != nil {
		return err
	}
	return fmt.Errorf(errorResp.Error)
}