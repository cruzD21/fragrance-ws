package db

import (
	"encoding/json"
	"fragrance-ws/models"
)

func (db *DatabaseConn) InsertAccordList(accordList []string, fragID int) error {

	for _, accord := range accordList {
		accordID, err := db.GetOrInsertAccord(models.Accord{
			Name:        accord,
			Description: "test accord description",
		})
		if err != nil {
			return err
		}

		err = db.InsertFragranceToAccord(accordID, fragID)
		if err != nil {
			return err
		}

	}
	return nil
}

func (db *DatabaseConn) GetOrInsertAccord(accord models.Accord) (int, error) {
	var err error
	var accordID int
	accordID, err = db.getNoteID(accord.Name)
	if err != nil {
		//note does not exist

		accordID, err = db.InsertAccord(accord)
		if err != nil {
			return 0, err
		}
	}

	return accordID, nil
}

func (db *DatabaseConn) InsertAccord(accord models.Accord) (int, error) {
	var res []struct {
		ID int `json:"accord_id"`
	}
	data, _, err := db.Supabase.
		From("accord").
		Insert(
			accord,
			false,
			"",
			"minimal",
			"exact").
		Execute()
	if err != nil {
		return 0, err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}

	return res[0].ID, err
}

func (db *DatabaseConn) getAccordID(accordName string) (int, error) {
	var err error
	var res []struct {
		ID int `json:"accord_id"`
	}
	data, _, err := db.Supabase.
		From("accord").
		Select("accord_id", "1", false).
		Eq("name", accordName).
		Execute()

	if err != nil {
		return 0, err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}

	return res[0].ID, err
}

func (db *DatabaseConn) InsertFragranceToAccord(accordID, fragID int) error {
	row := models.FragranceToAccord{
		FragID:   fragID,
		AccordID: accordID,
	}
	_, _, err := db.Supabase.
		From("fragrance_to_accord").
		Insert(
			row,
			false,
			"",
			"minimal",
			"exact").
		Execute()

	return err
}
