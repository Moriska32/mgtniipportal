package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

// Connecting to db
func Connect() *sql.DB {
	//"postgres://portaluser:PortalDB2020@172.20.0.82:5432/portaldb"
	db, err := sql.Open("postgres", "postgres://portaluser:PortalDB2020@172.20.0.82:5432/portaldb")

	if err != nil {
		log.Printf("Reason: %v\n", err)
		os.Exit(100)
	}
	log.Printf("Connected to db")
	return db
}
