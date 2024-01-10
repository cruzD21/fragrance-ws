package db

import (
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
	"os"
)

func databaseInit() (*supabase.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	apiURL := os.Getenv("DB_API_URL")
	apiKEY := os.Getenv("DB_API_KEY")
	client, err := supabase.NewClient(apiURL, apiKEY, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
