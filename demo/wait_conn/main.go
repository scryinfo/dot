// wait for the connection

package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {

	//const DevServerDbConf = `{"host":"127.0.0.1", "port":"5432", "user": "scry", "password": "scry", "database": "kim_nft_test", "showSql": false}`
	dsn := "postgres://scry:scry@127.0.0.1:5432/kim_nft_test?&timeout=5&sslmode=disable"

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	sqldb.SetMaxOpenConns(2)
	db := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	var txs []bun.Tx
	ctx := context.Background()
	for range 3 {
		// this will dead wait until the connect is valid
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			fmt.Println(err)
		} else {
			txs = append(txs, tx)
		}
	}
	for _, tx := range txs {
		tx.Commit()
	}

}
