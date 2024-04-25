package db

import (
	"encoding/json"
	"errors"
	"fragrance-ws/models"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
	"log"
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
	if fragrance.Name == "" {
		return 0, errors.New("fragrance name empty, do not insert")
	}
	var err error
	var res []struct {
		FragranceID int `json:"fragrance_id"`
	}

	fragrance.HouseID = houseID
	_, err = db.GetFragranceID(fragrance.Name)
	if err != nil {
		return 0, errors.New("fragrance already exists ib db, skipping")
	}

	log.Println("house id", houseID)
	data, _, err := db.Supabase.
		From("fragrance").
		Insert(
			fragrance,
			false,
			"",
			"representation",
			"exact",
		).Execute()

	if err != nil {
		return 0, err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}

	return res[0].FragranceID, nil
}

func (db *DatabaseConn) InsertIntoFragranceHouse(house models.FragranceHouse) (int, error) {
	if house.Name == "" {
		return 0, errors.New("empty name, do not insert")
	}
	var err error
	var res []struct {
		HouseID int `json:"house_id"`
	}
	houseID, err := db.GetHouseID(house.Name)
	if err == nil {
		return houseID, err // house already exists
	}
	data, _, err := db.Supabase.
		From("fragrance_house").
		Insert(
			house,
			false,
			"",
			"representation",
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

//func (db *DatabaseConn) InsertFailedURLS(url []string) {
//	var err error
//	if len(url) == 0 {
//		log.Println("no failed urls")
//		return
//	}
//
//	data, _, err := db.Supabase.
//		From("fragrance_house").
//		Insert(
//			house,
//			false,
//			"",
//			"representation",
//			"exact",
//		).Execute()
//
//	if err != nil {
//		return 0, err
//	}
//
//	if err = json.Unmarshal(data, &res); err != nil {
//		return 0, err
//	}
//
//	return res[0].HouseID, nil
//}

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

func (db *DatabaseConn) GetHouseID(HouseName string) (int, error) {
	if HouseName == "" {
		return 0, errors.New("empty house name return")
	}

	var err error
	var res []struct {
		ID int `json:"house_id"`
	}
	data, _, err := db.Supabase.
		From("fragrance_house").
		Select("house_id", "1", false).
		Eq("name", HouseName).
		Execute()

	if err != nil {
		return 0, err
	}

	// Check if the query returned any rows
	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}

	if len(res) == 0 {
		return 0, errors.New("no such house exists")
	}

	return res[0].ID, err
}
