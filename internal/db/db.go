package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"db_course_project/internal/config"
)

// Connect opens a sqlx DB handle with pgx and configures the pool.
func Connect(cfg config.DBConfig) (*sqlx.DB, error) {
	sqlDB, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	wrapped := sqlx.NewDb(sqlDB, "pgx")
	wrapped.MapperFunc(sqlx.NameMapper)

	// Verify connectivity on startup.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := wrapped.DB.PingContext(ctx); err != nil {
		return nil, err
	}

	return wrapped, nil
}
