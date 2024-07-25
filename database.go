package databasepzngo

import (
	"database/sql"
	"time"
)

func GetConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/pzn_go_db?parseTime=true") // parseTime=true biar return data date/time dari db langsung ke parse jadi tipe data time.Time di golang, kalo engga nanti harus konversi dari []uint8 ke string dan dari string diparse ke time.Time
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)  // close koneksi yg udah ga dipake
	db.SetConnMaxLifetime(60 * time.Minute) // buat renew koneksi yg udah lama

	return db
}
