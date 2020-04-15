package test

import (
	"context"

	. "allaboutapps.at/aw/go-mranftl-sample/models"
	"github.com/volatiletech/sqlboiler/boil"
)

// A common interface for all model instances do allow all of the to be exposed in a []Fixture
type Fixture interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
}

// Any colums non specified will be boil.Infer()ed
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

	language3 = Language{
		ID:       "20ecd47d-f453-4ab9-8e62-3289604211c0",
		Language: "SE",
	}

	pilotLanguage1 = PilotLanguage{
		LanguageID: language1.ID,
		PilotID:    pilot1.ID,
	}

	pilotLanguage2 = PilotLanguage{
		LanguageID: language1.ID,
		PilotID:    pilot2.ID,
	}

	pilotLanguage3 = PilotLanguage{
		LanguageID: language3.ID,
		PilotID:    pilot3.ID,
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

	// Defines the order in which everything will get inserted
	fixtures = []Fixture{
		&pilot1,
		&pilot2,
		&pilot3,
		&language1,
		&language2,
		&language3,
		&pilotLanguage1,
		&pilotLanguage2,
		&pilotLanguage3,
		&jet1,
		&jet2,
	}
)
