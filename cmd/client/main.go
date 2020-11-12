package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/mimatache/go-shop/pkg/client"
)

func main() {
	port := flag.String("port", "9090", "Port of server")

	flag.Parse()

	address := fmt.Sprintf("http://localhost:%s/api/v1", *port)
	clients := []*client.ShopClient{}
	john, err := client.New("john.doe@company.com", "1234", address)
	if err != nil {
		panic(err)
	}
	clients = append(clients, john)
	john2, err := client.New("john.doe2@company.com", "1234", address)
	if err != nil {
		panic(err)
	}
	clients = append(clients, john2)

	var wg sync.WaitGroup
	wg.Add(len(clients))
	for _, customer := range clients {
		go func(customer *client.ShopClient) {
			defer wg.Done()
			err := customer.PerformActionLoop(map[uint]uint{1: 2, 2: 1})
			if err != nil {
				fmt.Printf("user %s failed the buy: %v \n", customer.Name, err)
				return
			}
			fmt.Printf("user %s made the buy \n", customer.Name)
		}(customer)
	}
	wg.Wait()
}
