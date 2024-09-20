package migration

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunMigrations(pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	goose.SetBaseFS(nil)
	opts := goose.WithNoVersioning()
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations", opts); err != nil {
		panic(err)
	}
	if err := db.Close(); err != nil {
		panic(err)
	}
	return nil
}
