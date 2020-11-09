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

var (
	table = &productTable{name: "products"}
)

type productTable struct {
	name string
}

func (u *productTable) GetName() string {
	return u.name
}

func (u *productTable) GetTableSchema() *memdb.TableSchema {
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

func New(log logger, seed io.Reader) (ProductStore, error) {
	schema := store.NewSchema()

	var products []*Product

	err := json.NewDecoder(seed).Decode(&products)
	if err != nil {
		return nil, err
	}

	schema.AddToSchema(table)

	db, err := store.New(schema)
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		err = db.Write(table, product)
		if err != nil {
			return nil, err
		}
	}

	return &productLogger{
		log:  log,
		next: &productStore{db: db},
	}, nil
}

type ProductStore interface {
	GetProductById(ID uint) (*Product, error)
	SetProducts(products ...*Product) (store.ConditionMet, error)
}

type productStore struct {
	db store.Store
}

func (p *productStore) GetProductById(ID uint) (*Product, error) {
	raw, err := p.db.Read(table, "id", ID)
	return checkAndReturn(raw, err)
}

func (p *productStore) SetProducts(products ...*Product) (store.ConditionMet, error) {
	objs := make([]interface{}, len(products))
	for i, v := range products {
		objs[i] = v
	}
	return p.db.WriteAfterExternalCondition(table, objs...)
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

func (p *productLogger) GetProductById(ID uint) (*Product, error) {
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
	product, err = p.next.GetProductById(ID)
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
