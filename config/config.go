package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Connect to db
func Connect() *sql.DB {
	//"postgres://portaluser:PortalDB2020@172.20.0.82:5432/portaldb"
	db, err := sql.Open("postgres", "postgres://portaluser:PortalDB2020@172.20.0.82:5432/portaldb")

	if err != nil {
		log.Printf("Reason: %v\n", err)
		os.Exit(100)
	}

	return db
}
