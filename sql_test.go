package databasepzngo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) { // bisa dipake buat perintah non query data, insert/update/delete
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customer(id, name) VALUES('I', 'Izzuddin')"

	_, err := db.ExecContext(ctx, script)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name FROM customer"

	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // jangan lupa close rowsnya

	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name) // pake pointer buat nge-set data dari param
		if err != nil {
			panic(err)
		}
		fmt.Println("Id: ", id)
		fmt.Println("Name: ", name)
	}
}

func TestQuerySqlComplex(t *testing.T) { // gaada data yg null
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer" // direkomendasikan ditulis semua kolomnya, hindari * biar ga bingung positioning kalo ternyata ada alter table

	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // jangan lupa close rowsnya

	for rows.Next() {
		var id, name, email string
		var balance int32
		var rating float64
		var birthDate, createdAt time.Time
		var married bool
		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}
		fmt.Println("Id:", id, "Name:", name, "Email:", email, "Balance:", balance, "Rating:", rating, "Birth date:", birthDate, "Married:", married, "Created at:", createdAt)
	}
}

func TestQuerySqlComplexNullable(t *testing.T) { // ada data yg null
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"

	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // jangan lupa close rowsnya

	for rows.Next() {
		var id, name string
		var email sql.NullString // tipe datanya diganti jadi ini
		var balance int32
		var rating float64
		var birthDate sql.NullTime // ini juga diganti
		var createdAt time.Time
		var married bool
		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}
		fmt.Println("/------------------------------/")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)
		if email.Valid {
			fmt.Println("Email:", email.String)
		} else {
			fmt.Println("Email: NULL")
		}
		fmt.Println("Balance:", balance)
		fmt.Println("Rating:", rating)
		if birthDate.Valid {
			fmt.Println("Birth date:", birthDate.Time)
		} else {
			fmt.Println("Birth date: NULL")
		}
		fmt.Println("Married:", married)
		fmt.Println("Created at:", createdAt)
	}
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin'; #" // contoh sql injection
	password := "salah"

	script := "SELECT username FROM user WHERE username = '" + username + "' AND password = '" + password + "' LIMIT 1" // bisa kena sql injection karena sql querynya hardcoded

	//! tidak disarankan buat hardcoded kalo ngirim query ada parameternya

	// scriptnya jadi gini
	// SELECT username FROM user WHERE username = 'admin'; #' AND password = 'admin' LIMIT 1
	// setelah # querynya ga dianggep jadi cuman pake query username aja

	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // jangan lupa close rowsnya

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Success login", username)
	} else {
		fmt.Println("Login failed")
	}
}

func TestQuerySqlParams(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin'; #" // contoh sql injection, tapi bakal gagal
	password := "admin"

	script := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1" // gagal karena querynya pake params ?

	rows, err := db.QueryContext(ctx, script, username, password) // paramsnya disebutin di sini yang bakal gantiin ? secara urut
	if err != nil {
		panic(err)
	}
	defer rows.Close() // jangan lupa close rowsnya

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Success login", username)
	} else {
		fmt.Println("Login failed")
	}
}

func TestExecSqlParams(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	id := "H"
	name := "Hanif"

	script := "INSERT INTO customer(id, name) VALUES(?,?)"

	_, err := db.ExecContext(ctx, script, id, name)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "hanif@mail.com"
	comment := "Tes komen"

	script := "INSERT INTO comments(email, comment) VALUES(?,?)"

	result, err := db.ExecContext(ctx, script, email, comment)
	if err != nil {
		panic(err)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new comment with id", insertId)
}

func TestPrepareStatement(t *testing.T) { // dipake kalo input data banyak/berkali2 di query/exec yg sama
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO comments(email, comment) VALUES(?,?)"

	statement, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer statement.Close() // jangan lupa diclose

	for i := 0; i < 10; i++ {
		email := "luthfi" + strconv.Itoa(i) + "@mail.com"
		comment := "Tes komen ke-" + strconv.Itoa(i)

		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id: ", lastInsertId)
	}
}

func TestDatabaseTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO comments(email, comment) VALUES(?,?)"

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	// do transaction
	for i := 0; i < 10; i++ {
		email := "izzuddin" + strconv.Itoa(i) + "@mail.com"
		comment := "Tes komen ke-" + strconv.Itoa(i)

		result, err := tx.ExecContext(ctx, script, email, comment)
		if err != nil {
			panic(err)
		}

		lastInsertId, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id: ", lastInsertId)
	}

	err = tx.Commit() // bisa rollback kalo mau digagalin transaksinya
	if err != nil {
		panic(err)
	}
}
