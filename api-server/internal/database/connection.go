package database

import (
	"api-server/internal/configurator"
	"api-server/internal/database/migrations"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Database struct {
	Db *sql.DB
}

type DatabaseInterface interface {
	AddProduct(product Product) int
	GetProduct(id int) Product
	GetAllProducts() []Product
}

func CreateNewResourceRepository(ctx context.Context, config *configurator.SqlConfig) (DatabaseInterface, error) {
	sqlDb, err := CreateNewConnection(ctx, config)
	if err != nil {
		log.Fatal("Failed to create a new SQL connection")
		return nil, err
	}

	err = migrations.RunDatabaseMigrations(sqlDb, migrations.DefaultSourceUrl, migrations.CurrentDatabaseMigrationVersion)
	if err != nil {
		log.Fatal("Failed to run database migrations")
		return nil, err
	}

	return &Database{
		Db: sqlDb,
	}, nil
}

func CreateNewConnection(ctx context.Context, sqlConfig *configurator.SqlConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		sqlConfig.DbHost,
		sqlConfig.DbPort,
		sqlConfig.DbUsername,
		sqlConfig.DbPassword,
		sqlConfig.DbName,
		sqlConfig.DbSslMode,
	)

	sqlDb, err := sql.Open(sqlConfig.DbDriver, dsn)
	if err != nil {
		log.Fatal("Failed to open a connection with database")
		return nil, err
	}

	sqlDb.SetMaxIdleConns(sqlConfig.MaxIdleConns)
	sqlDb.SetMaxOpenConns(sqlConfig.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(sqlConfig.ConnMaxLifetime) * time.Minute)

	err = sqlDb.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
		return nil, err
	}

	return sqlDb, nil
}

func (db *Database) AddProduct(product Product) int {
	query := `
	INSERT INTO product (name, price, available)
	VALUES ($1, $2, $3) RETURNING id
	`

	var pk int
	err := db.Db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&pk)
	if err != nil {
		log.Fatal(err)
	}
	return pk
}

func (db *Database) GetProduct(id int) Product {
	query := `SELECT name, price, available FROM product WHERE id=$1`
	product := Product{}
	err := db.Db.QueryRow(query, id).Scan(&product.Name, &product.Price, &product.Available)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("No rows found with id %d", id)
		}
		log.Fatal(err)
	}
	fmt.Println(product)
	return product
}

func (db *Database) GetAllProducts() []Product {
	query := `SELECT name, price, available FROM product`
	product := Product{}
	products := []Product{}
	rows, err := db.Db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&product.Name, &product.Price, &product.Available)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, product)
	}
	fmt.Println(products)
	return products
}
