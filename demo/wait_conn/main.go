// wait for the connection

package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	if tx, err := db.BeginTx(context.Background(), nil); err == nil {
		fmt.Println("begin tx ok")
		txs = append(txs, tx)
	}
	if tx, err := db.BeginTx(context.Background(), nil); err == nil {
		fmt.Println("begin tx ok")
		txs = append(txs, tx)
	}

	// this will dead wait until the connect is valid
	ctxTimeout, c := context.WithTimeout(context.Background(), 5*time.Second)
	tx, err := db.BeginTx(ctxTimeout, nil)
	if err != nil {
		fmt.Println("begin tx error : ", err)
		fmt.Println("time out : ", ctxTimeout.Err())
	} else {
		select {
		case <-ctxTimeout.Done():
			fmt.Println("---- timeout no error --- ")
		default:
			fmt.Println("begin tx ok")
			txs = append(txs, tx)
		}
	}
	c()
	for _, tx := range txs {
		tx.Commit()
	}
	db.Close()
}
