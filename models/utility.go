package models

type NoteCategories struct {
	TopNotes    []string
	MiddleNotes []string
	BaseNotes   []string
}

type FragrancePage struct {
	Fragrance   Fragrance
	FragHouse   FragranceHouse
	NoteCat     NoteCategories
	MainAccords []string
}
