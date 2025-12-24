package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Try loading .env (optional). If present, environment variables will be populated.
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASEURL")
	}
	if dbURL == "" {
		log.Fatal("environment variable DATABASE_URL is required; set it in .env or export it and re-run")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// no connection lifetime/timeouts set; rely on defaults
	_, _ = db.DB()

	// include both singular and plural table names we've seen in this DB
	tables := []string{"patients", "staff", "staffs", "hospitals"}
	for _, t := range tables {
		qry := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", t)
		if err := db.Exec(qry).Error; err != nil {
			log.Fatalf("failed dropping table %s: %v", t, err)
		}
		fmt.Println("dropped table:", t)
	}

	fmt.Println("All specified tables dropped successfully.")

	// Confirmation: list remaining tables in the public schema
	fmt.Println("Verifying remaining tables in public schema:")
	rows, err := db.Raw("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public'").Rows()
	if err != nil {
		log.Fatalf("failed listing tables: %v", err)
	}
	defer rows.Close()
	var found bool
	for rows.Next() {
		var tbl string
		if err := rows.Scan(&tbl); err != nil {
			log.Fatalf("failed scanning table name: %v", err)
		}
		fmt.Println(" -", tbl)
		found = true
	}
	if !found {
		fmt.Println(" - <no tables found in public schema>")
	}
}
