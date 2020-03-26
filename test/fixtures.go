package test

import (
	"context"

	. "allaboutapps.at/aw/go-mranftl-sample/models"
	"github.com/volatiletech/sqlboiler/boil"
)

type Insertable interface {
	InsertP(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns)
}

func GetFixtures() ([]Pilot, []Jet, []Language, []PilotLanguage) {

	// anything not specified will be boil.Infer()ed

	pilots := []Pilot{
		{
			ID:   "0ed44d09-3d8e-43c2-b203-5da10e5e266f",
			Name: "Mario",
		},
		{
			ID:   "de5d41ca-62f9-4c22-acec-29e520807607",
			Name: "Nick",
		},
		{
			ID:   "69c591e7-14cb-4f33-bb46-653242f84034",
			Name: "Hi",
		},
	}

	languages := []Language{
		{
			ID:       "9513365a-6b2d-4b0c-9d89-fd9a5c891a68",
			Language: "DE",
		},
		{
			ID:       "3377ad0c-4386-433d-85eb-8f1ee1ea418a",
			Language: "EN",
		},
	}

	pilotLanguages := []PilotLanguage{
		{
			LanguageID: languages[0].ID,
			PilotID:    pilots[0].ID,
		},
		{
			LanguageID: languages[1].ID,
			PilotID:    pilots[1].ID,
		},
	}

	jets := []Jet{
		{
			ID:      "41495f88-4459-4949-8871-11aa5a8a2b72",
			PilotID: pilots[0].ID,
			Age:     34,
			Color:   "green",
			Name:    "Jet1",
		},
		{
			ID:      "3713821c-b7df-4f05-961f-6acd18051bba",
			PilotID: pilots[1].ID,
			Age:     23,
			Color:   "blue",
			Name:    "Jet2",
		},
	}

	return pilots, jets, languages, pilotLanguages
}
