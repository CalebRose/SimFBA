package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/managers"
	"github.com/joho/godotenv"
)

// We recreate the logic from main.go to ensure the DB is ready
func InitialMigration() {
	initiate := dbprovider.GetInstance().InitDatabase()
	if !initiate {
		log.Println("Initiate pool failure... Ending application")
		os.Exit(1)
	}
}

func main() {
	// 1. Load the .env from the parent directory
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file from parent directory")
		return
	}

	// 2. Run the actual migration/init logic from your main.go
	// This opens the SSH tunnel and sets up the GORM pool
	InitialMigration()

	fmt.Println("--- DATABASE INITIALIZED & TUNNEL OPEN ---")

	// 3. Run the logic
	managers.RunPreDraftEvents()

	fmt.Println("--- DRY RUN FINISHED ---")
}
