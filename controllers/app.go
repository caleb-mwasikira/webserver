package controllers

import (
	"database/sql"

	models "github.com/caleb-mwasikira/webserver/models"
)

type Application struct {
	// Repositories of where our application fetches data from
	// (currently fetching from the database)
	NotesRepo *models.NoteRepository
}

func NewApplication(db *sql.DB) *Application {
	return &Application{
		NotesRepo: &models.NoteRepository{
			DB: db,
		},
	}
}
