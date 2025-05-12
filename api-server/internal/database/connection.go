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
