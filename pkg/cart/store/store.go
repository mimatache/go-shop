package store

import (
	"github.com/hashicorp/go-memdb"

	"github.com/mimatache/go-shop/internal/store"
)

type logger interface {
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}

const id = "id"

var (
	table = &shoppingCartTable{name: "shoppingCart"}
)

type shoppingCartTable struct {
	name string
}

func (u *shoppingCartTable) GetName() string {
	return u.name
}

func (u *shoppingCartTable) GetTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: u.name,
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:    "id",
				Unique:  true,
				Indexer: &memdb.UintFieldIndex{Field: "ID"},
			},
			"products": &memdb.IndexSchema{
				Name:    "products",
				Unique:  false,
				Indexer: &memdb.StringMapFieldIndex{Field: "Products"},
			},
		},
	}
}

type CartStore interface {
	AddProductWhenConditionIsMet(userID uint, prodID uint, quantity uint) (store.ConditionMet, uint, error)
	GetProductsForUser(userID uint) (map[uint]uint, error)
	ClearCartFor(userID uint) error
}

func New(log logger) (CartStore, error) {
	schema := store.NewSchema()

	schema.AddToSchema(table)

	db, err := store.New(schema)
	if err != nil {
		return nil, err
	}

	return &cartLogger{
		log:  log,
		next: &cartStore{db: db},
	}, nil
}

type cartStore struct {
	db store.Store
}

func (c *cartStore) AddProductWhenConditionIsMet(userID uint, prodID uint, quantity uint) (store.ConditionMet, uint, error) {
	var cartItem *CartItem
	cartItem, err := c.getProductsForUser(userID)
	switch err.(type) {
	case nil:
		if val, ok := cartItem.Products[prodID]; ok {
			cartItem.Products[prodID] = val + quantity
		} else {
			cartItem.Products[prodID] = quantity
		}
	case store.NotFound:
		cartItem, err = NewCartItem(userID, prodID, quantity)
		if err != nil {
			return nil, 0, err
		}
	default:
		return nil, 0, err
	}
	isMet, err := c.db.WriteAfterExternalCondition(table, cartItem)
	return isMet, cartItem.Products[prodID], err
}

func (c *cartStore) getProductsForUser(userID uint) (*CartItem, error) {
	item, err := c.db.Read(table, id, userID)
	if err != nil {
		return nil, err
	}
	cartItem := item.(*CartItem)
	return cartItem, nil
}

func (c *cartStore) ClearCartFor(userID uint) error{
	return c.db.Remove(table, id, userID)
}

func (c *cartStore) GetProductsForUser(userID uint) (map[uint]uint, error) {
	cart, err := c.getProductsForUser(userID)
	if err != nil {
		return nil, err
	}
	return cart.Products, nil
}

type cartLogger struct {
	log  logger
	next CartStore
}

func (c *cartLogger) AddProductWhenConditionIsMet(userID uint, prodID uint, quantity uint) (store.ConditionMet, uint, error) {
	var err error
	defer func() {
		if err != nil {
			c.log.Debugf("could not update cart for user %d err: %s", userID, err.Error())
			return
		}
		c.log.Debugf("updated cart for user %d for product %d with %d", userID, prodID, quantity)
	}()
	conditionMet, quantity, err := c.next.AddProductWhenConditionIsMet(userID, prodID, quantity)
	return conditionMet, quantity, err
}

func (c *cartLogger) GetProductsForUser(userID uint) (map[uint]uint, error) {
	var err error
	var items map[uint]uint
	defer func() {
		if err != nil {
			c.log.Debugf("could not retrieve the products for user %d err: %s", userID, err.Error())
			return
		}
		c.log.Debugf("retrieved products for user %d", userID)
		c.log.Debugw("current cart for user", "user", userID, "products", items)
	}()

	items, err = c.next.GetProductsForUser(userID)
	return items, err
}

func (c *cartLogger) ClearCartFor(userID uint) error{
	var err error
	defer func() {
		if err != nil {
			c.log.Debugf("could not remove the cart for user %d err: %s", userID, err.Error())
			return
		}
		c.log.Debugf("cleared cart for user %d", userID)
	}()

	err = c.next.ClearCartFor(userID)
	return err
}
