package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kevinanielsen/go-fast-cdn/src/database"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

func main() {
	// Parse command line flags
	rollback := flag.Bool("rollback", false, "Rollback the media unification migration")
	flag.Parse()

	// Load executable path
	util.LoadExPath()

	// Connect to database
	database.ConnectToDB()

	// Execute migration or rollback
	var err error
	if *rollback {
		fmt.Println("Rolling back media unification migration...")
		err = database.RollbackMediaMigration()
		if err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		fmt.Println("Media unification migration rolled back successfully!")
	} else {
		fmt.Println("Running media unification migration...")
		err = database.RunMediaMigration()
		if err != nil {
			log.Fatalf("Failed to run migration: %v", err)
		}
		fmt.Println("Media unification migration completed successfully!")
	}

	os.Exit(0)
}
