package postgresql

import (
	"api/pkg/configuration"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect(cfg *configuration.Postgresql) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Name)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
