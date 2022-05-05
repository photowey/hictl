package database

import (
	`database/sql`
)

// Hictl the config of hictl cmd
type Hictl struct {
	DatabaseMap map[string]Config `json:"databases" yaml:"databases" toml:"databases"` // database map
}

func (h *Hictl) AcquireDatabase(database string) (Config, bool) {
	conf, ok := h.DatabaseMap[database]

	return conf, ok
}

type Config struct {
	Database string `json:"database" yaml:"database" toml:"database"` // database name -> tests
	Host     string `json:"host" yaml:"host" toml:"host"`             // database host -> 192.168.1.11
	Port     string `json:"port" yaml:"port" toml:"port"`             // database port -> 3307
	UserName string `json:"username" yaml:"username" toml:"username"` // database user name -> root
	Password string `json:"password" yaml:"password" toml:"password"` // database database password -> root
}

type Database struct {
	Name   string
	Tables []*Table `json:"tables"`
}

// Table database table data-model
type Table struct {
	Name    string         `json:"tableName"`
	Comment sql.NullString `json:"tableComment"`
	Columns []*Column      `json:"columns"`
	Indexs  []*Index       `json:"indexs"`
}

// Column database table column data-model
type Column struct {
	TableName     string         `json:"tableName"`
	ColumnName    string         `json:"columnName"`
	ColumnComment sql.NullString `json:"columnComment"`
	NotNull       string         `json:"notNull"`
	DataType      string         `json:"dataType"`
	DataLength    sql.NullString `json:"dataLength"`
	PrimaryKey    sql.NullString `json:"primaryKey"`
	MaxBit        sql.NullString `json:"maxBit"`
	MinBit        sql.NullString `json:"minBit"`
}

// Index database table index data-model
type Index struct {
	TableName   string `json:"tableName"`
	IndexName   string `json:"indexName"`
	IndexColumn string `json:"indexColumn"`
}
