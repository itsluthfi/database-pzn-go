package databasepzngo

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestOpenConnection(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/pzn_go_db")
	if err != nil {
		panic(err)
	}
	defer db.Close() // selalu close connection setelah selesai
}
