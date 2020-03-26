package test

import (
	"context"

	. "allaboutapps.at/aw/go-mranftl-sample/models"
	"github.com/volatiletech/sqlboiler/boil"
)

type Insertable interface {
	InsertP(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns)
}

// Any colums non specified will be boil.Infer()ed!
var (
	pilot1 = Pilot{
		ID:   "0ed44d09-3d8e-43c2-b203-5da10e5e266f",
		Name: "Mario",
	}

	pilot2 = Pilot{
		ID:   "de5d41ca-62f9-4c22-acec-29e520807607",
		Name: "Nick",
	}

	pilot3 = Pilot{
		ID:   "69c591e7-14cb-4f33-bb46-653242f84034",
		Name: "Hi",
	}

	language1 = Language{
		ID:       "9513365a-6b2d-4b0c-9d89-fd9a5c891a68",
		Language: "DE",
	}

	language2 = Language{
		ID:       "3377ad0c-4386-433d-85eb-8f1ee1ea418a",
		Language: "EN",
	}

	pilotLanguage1 = PilotLanguage{
		LanguageID: language1.ID,
		PilotID:    pilot1.ID,
	}

	pilotLanguage2 = PilotLanguage{
		LanguageID: language1.ID,
		PilotID:    pilot2.ID,
	}

	jet1 = Jet{
		ID:      "41495f88-4459-4949-8871-11aa5a8a2b72",
		PilotID: pilot1.ID,
		Age:     34,
		Color:   "green",
		Name:    "Jet1",
	}

	jet2 = Jet{
		ID:      "3713821c-b7df-4f05-961f-6acd18051bba",
		PilotID: pilot2.ID,
		Age:     23,
		Color:   "blue",
		Name:    "Jet2",
	}

	// This defines the order in which everything will get inserted
	fixtures = []Insertable{
		&pilot1,
		&pilot2,
		&pilot3,
		&language1,
		&language2,
		&pilotLanguage1,
		&pilotLanguage2,
		&jet1,
		&jet2,
	}
)
