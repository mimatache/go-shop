package store

import (
	"github.com/hashicorp/go-memdb"
)

type Table interface {
	GetName() string
	GetTableSchema() *memdb.TableSchema
}

type Schema interface {
	AddToSchema(table Table)
	initDB() (*memdb.MemDB, error)
}

type schema struct {
	schema *memdb.DBSchema
}

func (s *schema) AddToSchema(table Table) {
	s.schema.Tables[table.GetName()] = table.GetTableSchema()
}

func (s *schema) initDB() (*memdb.MemDB, error) {
	return memdb.NewMemDB(s.schema)
}
func NewSchema() Schema {
	tableShema := make(map[string]*memdb.TableSchema)
	return &schema{schema: &memdb.DBSchema{Tables: tableShema}}
}
