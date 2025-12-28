package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect(databaseURL string) error {
	var err error
	DB, err = sql.Open("mysql", databaseURL)
	if err != nil {
		return fmt.Errorf("error al abrir la conexión: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)

	log.Println("✓ Conexión a la base de datos establecida")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
