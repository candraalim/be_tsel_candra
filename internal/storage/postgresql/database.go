package postgresql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/candraalim/be_tsel_candra/config"
)

type Database struct {
	*sqlx.DB
	schema string
}

func NewDatabase(cfg *config.DatabaseConfig) *Database {
	if cfg == nil {
		panic("config is nil")
	}
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%v sslmode=disable search_path=%s",
			cfg.Username,
			cfg.Password,
			cfg.Name,
			cfg.Host,
			cfg.Port,
			cfg.Schema))

	if err != nil {
		fmt.Println("failed to connect database")
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("failed to ping database")
		panic(err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenConn)

	return &Database{
		db,
		cfg.Schema,
	}
}

func (d *Database) SchemaName() string {
	return d.schema
}
