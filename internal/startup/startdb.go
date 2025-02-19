package startup

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Username string `toml:"username"`
}

func readConfig() Config {
	config := Config{}
	if _, err := toml.DecodeFile("settings.toml", &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func StartDB() {
	config := readConfig()
	username := config.Username

	if username == "" {
		log.Fatal("username is empty in settings.toml")
	}

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Build an absolute path for the database file
	dataDir := filepath.Join(wd, "data")
	dbPath := filepath.Join(dataDir, username+".db")

	// Ensure the data directory exists
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			log.Fatalf("failed to create data directory: %v", err)
		}
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// Create the database file if it does not exist
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.Printf("Database created at %s", dbPath)
	} else {
		log.Printf("Database already exists at %s", dbPath)
	}
}
