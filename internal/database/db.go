package database

import (
	"database/sql"
	"fmt"
	"go_code/simplek8s/server"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() *sql.DB {
	var db *sql.DB
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				server.Logger.Info("Database connection established")
				return db
			}
		}

		server.Logger.Error(fmt.Sprintf("Failed to connect to database, retrying in 5 seconds... (attempt %d/10)", i+1))
		time.Sleep(5 * time.Second)
	}
	server.Logger.Fatal(fmt.Sprintf("Could not connect to the database: %v", err))
	return nil
}
