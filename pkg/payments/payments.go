package payments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"crypto/rand"
)

type Luke struct {
	Name string `json:"name"`
}

func New() *API {
	return &API{
		client: &http.Client{Timeout: 2 * time.Second},
	}
}

type API struct {
	client *http.Client
}

// MakePayment randomly tells you that you can't pay for things
func (a *API) MakePayment(user string, money uint) error {
	time.Sleep(2 * time.Second)
	r, err := a.client.Get("https://swapi.dev/api/people/1/")
	if err != nil {
		return err
	}
	luke := &Luke{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(luke); err != nil {
		luke.Name = "Darth Vader"
	}
	if (randomNumber() % 2) == 0 {
		return fmt.Errorf("request blocked by %s", luke.Name)
	}
	return nil
}

func randomNumber() int {
	b := make([]byte, 1)
	_, err := rand.Read(b) 
	if err != nil{
		return 1
	}
	return int(b[0])
}
