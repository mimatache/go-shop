package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/mimatache/go-shop/internal/http/authorization"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log, flush, err := logger.New("shop", true)
	if err != nil {
		fmt.Printf("Could not instantiate logger %v", err)
		os.Exit(1)
	}
	defer flush()

	authorization.CleanBlacklist(time.Duration(5) * time.Second)
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
		log.Fatal(fmt.Sprintf("could not start DB %v", err))
	}

	// Loading the seeds to the DB
	err = userStore.LoadSeeds(userSeeds, db)
	if err != nil {
		log.Fatal(fmt.Sprintf("could not load seeds for user to DB %v", err))
	}
	err = productsStore.LoadSeeds(productSeeds, db)
	if err != nil {
		log.Fatal(fmt.Sprintf("could not load seeds for product to DB %v", err))
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

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", *port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-c
	log.Info("Shutting down...")
	ctx, cancel = context.WithTimeout(ctx, time.Second * 5)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
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
