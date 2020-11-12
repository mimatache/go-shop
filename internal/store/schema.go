package store

import (
	"github.com/hashicorp/go-memdb"
)

// Named represents a named object
type Named interface {
	GetName() string
}

// Table represents the name and schema of a DB
type Table interface {
	Named
	GetTableSchema() *memdb.TableSchema
}

// Schema represents the db schema
type Schema interface {
	AddToSchema(table Table)
	initDB() (*memdb.MemDB, error)
}

type schema struct {
	schema *memdb.DBSchema
}

// AddToSchema add a table schema to the schema
func (s *schema) AddToSchema(table Table) {
	s.schema.Tables[table.GetName()] = table.GetTableSchema()
}

func (s *schema) initDB() (*memdb.MemDB, error) {
	return memdb.NewMemDB(s.schema)
}

// NewSchema start a new DB schema
func NewSchema() Schema {
	tableShema := make(map[string]*memdb.TableSchema)
	return &schema{schema: &memdb.DBSchema{Tables: tableShema}}
}
