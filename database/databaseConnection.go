package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var PG_Client *sql.DB = DBinstance()

func DBinstance() *sql.DB {
	username := os.Getenv("PGUSER")
	if username == "" {
		log.Fatal("Username not found")
	}

	password := os.Getenv("PGPASSWORD")
	if password == "" {
		log.Fatal("Password not found")
	}

	PgDB := os.Getenv("PG_URL")
	if PgDB == "" {
		PgDB = "postgres://" + username + ":" + password + "@localhost:5432/go-oauth?sslmode=disable"
	}

	db, err := sql.Open("postgres", PgDB)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Postgresql")

	return db
}

func CreateUserTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		token TEXT,
		refresh_token TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		user_id UUID DEFAULT gen_random_uuid()
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}
}

func EnablePgCrypto(db *sql.DB) {
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	if err != nil {
		log.Fatalf("Error enabling pgcrypto extension: %v", err)
	}
}