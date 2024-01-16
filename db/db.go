package db

import (
	"encoding/json"
	"fragrance-ws/models"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
	"os"
)

type DatabaseConn struct {
	Supabase *supabase.Client
}

func (db *DatabaseConn) DatabaseInit() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	apiURL := os.Getenv("DB_API_URL")
	apiKEY := os.Getenv("DB_API_KEY")
	client, err := supabase.NewClient(apiURL, apiKEY, nil)
	if err != nil {
		return err
	}
	db.Supabase = client
	return nil
}

func (db *DatabaseConn) InsertPage(page models.FragrancePage) error {
	var houseID, fragID int
	var err error

	//insert house
	if houseID, err = db.InsertIntoFragranceHouse(page.FragHouse); err != nil {
		return err
	}

	//insert fragrance
	if fragID, err = db.InsertIntoFragrances(page.Fragrance, houseID); err != nil {
		return err
	}
	//insert note
	if err = db.InsertNoteCategories(page.NoteCat, fragID); err != nil {
		return err
	}
	//insert relationships
	if err = db.InsertAccordList(page.MainAccords, fragID); err != nil {
		return err
	}
	//&& a to f

	return nil
}

func (db *DatabaseConn) InsertIntoFragrances(fragrance models.Fragrance, houseID int) (int, error) {
	var err error
	var res []struct {
		FragranceID int `json:"fragrance_id"`
	}

	fragrance.HouseID = houseID
	data, _, err := db.Supabase.
		From("fragrance").
		Insert(
			fragrance,
			false,
			"",
			"representation",
			"exact",
		).
		Execute()

	if err != nil {
		return 0, err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}

	return res[0].FragranceID, nil
}

func (db *DatabaseConn) InsertIntoFragranceHouse(house models.FragranceHouse) (int, error) {
	var err error
	var res []struct {
		HouseID int `json:"house_id"`
	}
	data, _, err := db.Supabase.
		From("fragrance_house").
		Insert(
			house,
			false,
			"",
			"minimal",
			"exact",
		).Execute()
	if err != nil {
		return 0, err
	}

	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	return res[0].HouseID, nil
}

func (db *DatabaseConn) GetFragranceID(fragName string) (int, error) {
	var err error
	var res []struct {
		ID int `json:"fragrance_id"`
	}
	data, _, err := db.Supabase.
		From("fragrance").
		Select("fragrance_id", "1", false).
		Eq("name", fragName).
		Execute()

	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}
	return res[0].ID, err
}
