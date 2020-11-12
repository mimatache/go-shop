package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/http/middleware"
	"github.com/mimatache/go-shop/internal/logger"
	"github.com/mimatache/go-shop/internal/store"
	"github.com/mimatache/go-shop/pkg/cart"
	cartStore "github.com/mimatache/go-shop/pkg/cart/store"
	"github.com/mimatache/go-shop/pkg/payments"
	"github.com/mimatache/go-shop/pkg/products"
	productsStore "github.com/mimatache/go-shop/pkg/products/store"
	"github.com/mimatache/go-shop/pkg/users"
	userStore "github.com/mimatache/go-shop/pkg/users/store"
)

var (
	userSeeds    *os.File
	productSeeds *os.File
	port         *string
)

func main() {

	log, flush, err := logger.New("shop", true)
	if err != nil {
		fmt.Printf("Could not instantiate logger %v", err)
		os.Exit(1)
	}
	defer flush()

	// Reading seed files
	readFlagValues(log)

	log.Info("Starting app")

	r := mux.NewRouter()
	r.Use(middleware.Logging(log))
	versionedRouter := r.PathPrefix("/api/v1").Subrouter()

	// Starting DB instance
	schema := store.NewSchema()
	schema.AddToSchema(userStore.GetTable())
	schema.AddToSchema(productsStore.GetTable())
	schema.AddToSchema(cartStore.GetTable())
	db, err := store.New(schema)
	if err != nil {
		log.Errorf("could not start DB %v", err)
		os.Exit(1)
	}

	// Loading the seeds to the DB
	err = userStore.LoadSeeds(userSeeds, db)
	if err != nil {
		log.Errorf("could not load seeds for user to DB %v", err)
		os.Exit(1)
	}
	err = productsStore.LoadSeeds(productSeeds, db)
	if err != nil {
		log.Errorf("could not load seeds for product to DB %v", err)
		os.Exit(1)
	}

	// Starting user API
	userLogger := logger.WithFields(log, map[string]interface{}{"api": "users"})
	users.NewAPI(userLogger, versionedRouter, db)

	// Starting product API
	productLogger := logger.WithFields(log, map[string]interface{}{"api": "products"})
	productsAPI := products.NewAPI(productLogger, db)

	// Starting cart API
	cartLogger := logger.WithFields(log, map[string]interface{}{"api": "cart"})
	cart.NewAPI(cartLogger, productsAPI, payments.New(), db, versionedRouter, middleware.JWTAuthorization)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", *port), r); err != nil {
		log.Error(err)
	}
}

func readFlagValues(log logger.Logger) {
	var err error
	port = flag.String("port", "9090", "Port of server")
	userSeedsFile := flag.String("users", "data/users.json", "seed users to store")
	productSeedsFile := flag.String("products", "data/products.json", "seed products to store")
	flag.Parse()

	log.Infof("Reading user seed file: %s", *userSeedsFile)
	userSeeds, err = os.Open(*userSeedsFile)
	if err != nil {
		log.Errorf("could not read contents of seed file: %v", err)
		os.Exit(1)
	}

	log.Infof("Reading product seed file: %s", *productSeedsFile)
	productSeeds, err = os.Open(*productSeedsFile)
	if err != nil {
		log.Errorf("could not read contents of seed file: %v", err)
		os.Exit(1)
	}
}
