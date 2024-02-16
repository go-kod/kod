package ksql

import (
	"database/sql"
	"time"

	"dario.cat/mergo"
	"github.com/samber/lo"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
)

type Config struct {
	DriverName      string
	Dsn             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func (c Config) Build() *sql.DB {
	lo.Must0(mergo.Merge(&c, Config{
		DriverName:      "sqlite3",
		Dsn:             "",
		MaxIdleConns:    25,
		MaxOpenConns:    500,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: time.Hour,
	}))

	if c.Dsn == "" {
		panic("sql: DSN is empty")
	}

	db := lo.Must(otelsql.Open(c.DriverName, c.Dsn,
		otelsql.WithDBSystem(c.DriverName),
	))

	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)
	db.SetConnMaxIdleTime(c.ConnMaxIdleTime)

	return db
}
