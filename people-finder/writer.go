package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	w "stillrunning.org/people-finder/wikidata"
)

// Function to initialize the database and create the table
func initDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./people.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	// Create the table with unique index on ID and compound index on Age and SiteLinksCnt
	query := `
	CREATE TABLE IF NOT EXISTS persons (
		id TEXT PRIMARY KEY,
		name TEXT,
		birthDate TEXT,
		deathDate TEXT,
		pic TEXT,
		siteLinksCnt INTEGER,
		age INTEGER
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_person_id ON persons(id);
	CREATE INDEX IF NOT EXISTS idx_age_sitelinkcnt ON persons(age, siteLinksCnt);
	`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Database and tables initialized successfully.")

	return db
}

// Function to insert or update a person into the database (upsert)
func upsertPerson(db *sql.DB, p w.Person) error {
	query := `
	INSERT OR REPLACE INTO persons (id, name, birthDate, deathDate, pic, siteLinksCnt, age)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, p.Id, p.Name, p.BirthDate, p.DeathDate, p.Pic, p.SiteLinksCnt, p.Age)
	return err
}
