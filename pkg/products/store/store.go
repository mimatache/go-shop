package store

import (
	"encoding/json"
	"io"

	"github.com/hashicorp/go-memdb"

	"github.com/mimatache/go-shop/internal/store"
)

type logger interface {
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}

// UnderlyingStore represents the interface that the DB should implement to be usable
type UnderlyingStore interface {
	Read(table string, key string, value interface{}) (interface{}, error)
	WriteAfterExternalCondition(table string, objs ...interface{}) (store.ConditionMet, error)
	Write(table string, value ...interface{}) error
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
			"id": &memdb.IndexSchema{
				Name:    "id",
				Unique:  true,
				Indexer: &memdb.UintFieldIndex{Field: "ID"},
			},
			"name": &memdb.IndexSchema{
				Name:    "name",
				Unique:  false,
				Indexer: &memdb.StringFieldIndex{Field: "Name"},
			},
			"price": &memdb.IndexSchema{
				Name:    "price",
				Unique:  false,
				Indexer: &memdb.UintFieldIndex{Field: "Price"},
			},
			"stock": &memdb.IndexSchema{
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

type ProductStore interface {
	GetProductByID(ID uint) (*Product, error)
	SetProducts(products ...*Product) (store.ConditionMet, error)
}

type productStore struct {
	db UnderlyingStore
}

func (p *productStore) GetProductByID(ID uint) (*Product, error) {
	raw, err := p.db.Read(table.GetName(), "id", ID)
	return checkAndReturn(raw, err)
}

func (p *productStore) SetProducts(products ...*Product) (store.ConditionMet, error) {
	objs := make([]interface{}, len(products))
	for i, v := range products {
		objs[i] = v
	}
	return p.db.WriteAfterExternalCondition(table.GetName(), objs...)
}

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

func (p *productLogger) GetProductByID(ID uint) (*Product, error) {
	var err error
	var product *Product
	defer func() {
		if err != nil {
			p.log.Debugw("error occured when retrieving product", "id", ID, "err", err)
			return
		}
		p.log.Debugf("Retrieved product %d", ID)
		p.log.Debugw("Current stock for item", "id", ID, "stock", product)

	}()
	product, err = p.next.GetProductByID(ID)
	return product, err
}

func (p *productLogger) SetProducts(products ...*Product) (store.ConditionMet, error) {
	var err error
	defer func() {
		if err != nil {
			p.log.Debugw("error occured when setting products", "error", err)
			return
		}
		p.log.Debugf("Products successfully updated")
	}()
	condition, err := p.next.SetProducts(products...)
	return condition, err
}
