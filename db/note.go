package db

import (
	"encoding/json"
	"fragrance-ws/models"
	"log"
)

func (db *DatabaseConn) InsertNoteCategories(notes models.NoteCategories, fragID int) error {
	var err error

	if err = db.processNoteCategory("top_note", fragID, notes.TopNotes); err != nil {
		return err
	}
	if err = db.processNoteCategory("middle_note", fragID, notes.MiddleNotes); err != nil {
		return err
	}
	if err = db.processNoteCategory("base_note", fragID, notes.BaseNotes); err != nil {
		return err
	}

	return nil
}

func (db *DatabaseConn) GetOrInsertNote(note models.Note) (int, error) {
	var err error
	var noteID int
	noteID, err = db.getNoteID(note.Name)
	if err != nil {

		noteID, err = db.InsertNote(note)
		if err != nil {
			//note exists
			return 0, err
		}
	}

	return noteID, nil
}

func (db *DatabaseConn) processNoteCategory(noteType string, fragID int, noteList []string) error {

	for _, noteName := range noteList {
		noteID, err := db.GetOrInsertNote(models.Note{
			Name:        noteName,
			Description: "test description",
		})
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

	if err != nil {
		return 0, err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}

	return res[0].ID, nil
}

func (db *DatabaseConn) InsertNote(note models.Note) (int, error) {
	var err error
	var res []struct {
		NoteID int `json:"note_id"`
	}
	log.Println(note)
	data, _, err := db.Supabase.
		From("note").
		Insert(
			note,
			false,
			"",
			"representation",
			"exact").
		Execute()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	if err = json.Unmarshal(data, &res); err != nil {
		return 0, err
	}

	return res[0].NoteID, nil
}

func (db *DatabaseConn) InsertFragranceToNote(relationship string, fragID int, noteID int) error {
	row := models.FragranceToNote{
		FragID:   fragID,
		NoteID:   noteID,
		NoteType: relationship,
	}

	_, _, err := db.Supabase.
		From("fragrance_to_note").
		Insert(
			row,
			false,
			"",
			"minimal",
			"exact").
		Execute()

	return err
}
