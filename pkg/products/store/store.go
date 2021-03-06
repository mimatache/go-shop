package store

import (
	"encoding/json"
	"io"

	"github.com/hashicorp/go-memdb"

	"github.com/mimatache/go-shop/internal/store"
)

//go:generate mockgen -source ./store.go -destination mocks/store.go

// ProductTransaction is used to externalize transactions
type ProductTransaction struct {
	tableName   string
	transaction store.Transaction
}

// Commit is used to finalyze this transaction
func (r *ProductTransaction) Commit() {
	r.transaction.Commit()
}

// Abort is used to cancel this transaction
func (r *ProductTransaction) Abort() {
	r.transaction.Abort()
}

// Write inserts intp storace
func (r *ProductTransaction) Write(product *Product) error {
	return r.transaction.Insert(r.tableName, product)
}

type logger interface {
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}

// UnderlyingStore represents the interface that the DB should implement to be usable
type UnderlyingStore interface {
	Read(table string, key string, value interface{}) (interface{}, error)
	Write(table string, value ...interface{}) error
	WriteAndBlock(table string, value ...interface{}) (store.Transaction, error)
}

var (
	table = &ProductTable{name: "products"}
)

// GetTable returns the product schema
func GetTable() *ProductTable {
	return table
}

// LoadSeeds write the seed information to the DB.
func LoadSeeds(seed io.Reader, db UnderlyingStore) error {
	var products []*Product

	err := json.NewDecoder(seed).Decode(&products)
	if err != nil {
		return err
	}

	for _, product := range products {
		err = db.Write(table.GetName(), product)
		if err != nil {
			return err
		}
	}
	return nil
}

// ProductTable represents the product table in the DB
type ProductTable struct {
	name string
}

// GetName return the name of the product table
func (u *ProductTable) GetName() string {
	return u.name
}

// GetTableSchema returns the schema for the products table
func (u *ProductTable) GetTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: u.name,
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:    "id",
				Unique:  true,
				Indexer: &memdb.UintFieldIndex{Field: "ID"},
			},
			"name": {
				Name:    "name",
				Unique:  false,
				Indexer: &memdb.StringFieldIndex{Field: "Name"},
			},
			"price": {
				Name:    "price",
				Unique:  false,
				Indexer: &memdb.UintFieldIndex{Field: "Price"},
			},
			"stock": {
				Name:    "stock",
				Unique:  true,
				Indexer: &memdb.UintFieldIndex{Field: "Stock"},
			},
		},
	}
}

// New returns a new instance of ProductStore
func New(log logger, db UnderlyingStore) ProductStore {
	return &productLogger{
		log: log,
		next: &productStore{
			db: db,
		},
	}
}

// ProductStore models the Product DB
type ProductStore interface {
	GetProductByID(ID uint) (*Product, error)
	SetProducts(products ...*Product) (*ProductTransaction, error)
}

type productStore struct {
	db UnderlyingStore
}

// GetProductByID returns a product give the product ID
func (p *productStore) GetProductByID(id uint) (*Product, error) {
	raw, err := p.db.Read(table.GetName(), "id", id)
	return checkAndReturn(raw, err)
}

// SetProducts updates the Product DB with the given products
func (p *productStore) SetProducts(products ...*Product) (*ProductTransaction, error) {
	objs := make([]interface{}, len(products))
	for i, v := range products {
		objs[i] = v
	}
	transaction, err := p.db.WriteAndBlock(table.GetName(), objs...)
	if err != nil {
		return nil, err
	}
	return &ProductTransaction{
		tableName:   table.GetName(),
		transaction: transaction,
	}, nil
}

// checkAndReturn reads the output from the DB and returns a Product instance if no error occurred.
// This will panic if the DB does not return and error but the output is not an Product.
// Intentinally left to do this as if this happens it means we have an incosistency in the DB that should be resolve immediately
// and silent handling might mask this issue
func checkAndReturn(raw interface{}, err error) (*Product, error) {
	if err != nil {
		return nil, err
	}
	return raw.(*Product), nil
}

type productLogger struct {
	log  logger
	next ProductStore
}

func (p *productLogger) GetProductByID(id uint) (*Product, error) {
	var err error
	var product *Product
	defer func() {
		if err != nil {
			p.log.Debugw("error occurred when retrieving product", "id", id, "err", err)
			return
		}
		p.log.Debugf("Retrieved product %d", id)
		p.log.Debugw("Current stock for item", "id", id, "stock", product)
	}()
	product, err = p.next.GetProductByID(id)
	return product, err
}

func (p *productLogger) SetProducts(products ...*Product) (*ProductTransaction, error) {
	var err error
	defer func() {
		if err != nil {
			p.log.Debugw("error occurred when setting products", "error", err)
			return
		}
		p.log.Debugf("Products successfully updated")
	}()
	transaction, err := p.next.SetProducts(products...)
	return transaction, err
}
