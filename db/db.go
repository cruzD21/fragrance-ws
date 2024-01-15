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

func (db *DatabaseConn) InsertNoteCategories(notes models.NoteCategories, fragID int) error {
	var err error

	if err = db.processNoteCategory("top_note", fragID, notes.TopNotes); err != nil {
		return err
	}
	if err = db.processNoteCategory("middle_note", fragID, notes.MiddleNotes); err != nil {
		return err
	}
	if err = db.processNoteCategory("base_note", fragID, notes.TopNotes); err != nil {
		return err
	}

	return nil
}

func (db *DatabaseConn) processNoteCategory(noteType string, fragID int, noteList []string) error {

	for _, noteName := range noteList {
		noteID, err := db.GetOrInsertNote(noteName)
		if err != nil {
			return err
		}
		err = db.InsertFragranceToNote(noteType, fragID, noteID)
		if err != nil {
			return err
		}
	}
	return nil
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

func (db *DatabaseConn) getNoteID(noteName string) (int, error) {
	var err error
	var res []struct {
		ID int `json:"note_id"`
	}
	data, _, err := db.Supabase.
		From("note").
		Select("note_id", "1", false).
		Eq("name", noteName).
		Execute()

	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}

	return res[0].ID, err
}

func (db *DatabaseConn) GetOrInsertNote(noteName string) (int, error) {
	var err error
	var noteID int
	noteID, err = db.getNoteID(noteName)
	if err != nil {
		//note does not exist
		row := models.Note{
			Name:        noteName,
			Description: "test note",
		}
		noteID, err = db.InsertNote(row)
		if err != nil {
			return 0, err
		}
	}

	return noteID, nil
}

func (db *DatabaseConn) InsertNote(note models.Note) (int, error) {
	var err error
	var res []struct {
		NoteID int `json:"note_id"`
	}
	data, _, err := db.Supabase.
		From("note").
		Insert(
			note,
			false,
			"",
			"representation",
			"exact").
		Execute()

	if err = json.Unmarshal(data, &res); err != nil {
		//handle error
		_ = err
	}

	return res[0].NoteID, nil
}

func (db *DatabaseConn) InsertFragranceToNote(relationship string, fragID int, noteID int) error {
	row := models.FragranceToNote{
		FragID:   fragID,
		NoteID:   noteID,
		NoteType: relationship,
	}

	_, _, err := db.Supabase.From("fragrance_to_note").
		Insert(
			row,
			false,
			"",
			"minimal",
			"exact").
		Execute()

	return err
}

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
