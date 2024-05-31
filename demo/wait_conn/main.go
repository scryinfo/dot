// wait for the connection

package main

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {

	//const DevServerDbConf = `{"host":"127.0.0.1", "port":"5432", "user": "scry", "password": "scry", "database": "kim_nft_test", "showSql": false}`
	dsn := "postgres://scry:scry:127.0.0.1:5432:kim_nft_test?&timeout=5&sslmode=disable"

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	var done bool

	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	//todo something

	done = true
	tx.Commit()

}
