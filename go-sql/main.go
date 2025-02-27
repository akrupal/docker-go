// 1) docker pull postgres
// 2) docker run --name pg-container -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres
// 3) docker ps (to see the list of images)
// 4) docker exec -ti pg-container createdb -U postgres gopgtest
// 5) docker exec -ti pg-container psql -U postgres
// 	in the command that comes type:\c gopgtest
// 6) this created a database by the name gopgtest
// 7) \q to quit or exit out of database
// 8) run the below code it will perform modifications like create table insert data
// 9) to check the modifications run 5 followed by \dt (dt is display table in short)
// 10) to actually check the data inside table you can run sql queries
// 		SELECT * FROM product;

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Product struct {
	Name      string
	Price     float64
	Available bool
}

func main() {
	connstr := "postgres://postgres:secret@localhost:5432/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	createProductTable(db)
	p := Product{"book", 12.88, true}

	pk := insertProduct(db, p)
	// fmt.Println("id of book is", pk)
	var name string
	var price float64
	var available bool

	query := `SELECT name, price, available FROM product WHERE id=$1`

	err = db.QueryRow(query, pk).Scan(&name, &price, &available)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No rows found with ID %d", pk)
			//this error will never happen in current case but lets say we have some product key and try it
			// then this error will be useful
		}
		log.Fatal(err)
	}
	fmt.Println("name is", name)
	fmt.Println("price is", price)
	fmt.Println("availability is", available)

	data := []Product{}
	rows, err := db.Query("SELECT name, price, available FROM product")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &price, &available)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Product{name, price, available})
	}
	fmt.Println(data)
}

func createProductTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS product (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(6,2) NOT NULL,
		available BOOLEAN,
		created TIMESTAMP DEFAULT NOW()
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

}

func insertProduct(db *sql.DB, product Product) int {
	query := `
	INSERT INTO product (name, price, available)
	VALUES ($1, $2, $3) RETURNING id
	`
	var pk int
	err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&pk)
	if err != nil {
		log.Fatal(err)
	}
	return pk
}
