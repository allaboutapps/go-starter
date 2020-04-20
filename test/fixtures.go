package test

import (
	"context"

	. "allaboutapps.at/aw/go-mranftl-sample/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

// A common interface for all model instances do allow all of the to be exposed in a []Fixture
type Fixture interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
}

// Any colums non specified will be boil.Infer()ed
var (
	user1 = User{
		ID:       "f6ede5d8-e22a-4ca5-aa12-67821865a3e5",
		IsActive: true,
		Username: null.StringFrom("user1@example.com"),
		Password: null.StringFrom("d79f8266fc2942c27af782e2daab181fb6948835fbef33ec4d6c1f3886b3d901381c048a4539bf4ae2ae01d31016e1ca7432917dabd11db27dbc28579dfbfee007f685746788e2bea36ff07fdc45806ff48e9551634b2675534d1169114d5564443c350476df985d6be0cec0fea0d7e3b248c7022a0d9cecbd10d481aecf7b0c97c5756a4f640138fa459729e0921e05c77e77c0a7b025007f2c6d780a1824fe404030f3b9117130818deb370465a9f5f9e32d5508ab06cb8867283b524c6e3a02e7002c95a69fca95f789f82c53d8ab3b2e435983ab25a8544ee9db3cd896c1fc2e9ba24c528754665a0785fe4f25d9d72e47807804d09360e38dc11d02febad845dd233c5227f8cd8c970ab3b117ce9d000b6af0aaf01a6ee6d82627e022ece143642b157f5f883e5fd930deaf4a253e5c5566fbd9ccbd2e9d8a3605dfe979856ca3e743b8d3027448718bba21d73baeb2d459f909b66f257c01bd76e23bf09da8ae1bbab5a6c6cef8b23ec260f25db4d9191a14142db762f4043896c49f6aca4b317f00394715612107b4b360d7120fda3022a7d67aa11380acaa676e43195f527824cf12346898040ed77bb6eb760c3b5861777809b6cfd0e57257584803e18ef580c07eeb0102f04603676702cc83c26aa30bc2bdca2cea51fe5937ff3616467ab9f40a8f1192669f06a387d2d9ed6831fd6371f4c6ad309e9730482a07"),
		Salt:     null.StringFrom("3f9dd9e6d89c7b545ad2f22cf4f68f25188fbf1b2c2fec30c2884c939eead8b5"),
	}

	// pilot1 = Pilot{
	// 	ID:   "0ed44d09-3d8e-43c2-b203-5da10e5e266f",
	// 	Name: "Mario",
	// }

	// pilot2 = Pilot{
	// 	ID:   "de5d41ca-62f9-4c22-acec-29e520807607",
	// 	Name: "Nick",
	// }

	// pilot3 = Pilot{
	// 	ID:   "69c591e7-14cb-4f33-bb46-653242f84034",
	// 	Name: "Hi",
	// }

	// language1 = Language{
	// 	ID:       "9513365a-6b2d-4b0c-9d89-fd9a5c891a68",
	// 	Language: "DE",
	// }

	// language2 = Language{
	// 	ID:       "3377ad0c-4386-433d-85eb-8f1ee1ea418a",
	// 	Language: "EN",
	// }

	// language3 = Language{
	// 	ID:       "20ecd47d-f453-4ab9-8e62-3289604211c0",
	// 	Language: "SE",
	// }

	// pilotLanguage1 = PilotLanguage{
	// 	LanguageID: language1.ID,
	// 	PilotID:    pilot1.ID,
	// }

	// pilotLanguage2 = PilotLanguage{
	// 	LanguageID: language1.ID,
	// 	PilotID:    pilot2.ID,
	// }

	// pilotLanguage3 = PilotLanguage{
	// 	LanguageID: language3.ID,
	// 	PilotID:    pilot3.ID,
	// }

	// jet1 = Jet{
	// 	ID:      "41495f88-4459-4949-8871-11aa5a8a2b72",
	// 	PilotID: pilot1.ID,
	// 	Age:     34,
	// 	Color:   "green",
	// 	Name:    "Jet1",
	// }

	// jet2 = Jet{
	// 	ID:      "3713821c-b7df-4f05-961f-6acd18051bba",
	// 	PilotID: pilot2.ID,
	// 	Age:     23,
	// 	Color:   "blue",
	// 	Name:    "Jet2",
	// }

	// Defines the order in which everything will get inserted
	fixtures = []Fixture{
		&user1,
		// &pilot1,
		// &pilot2,
		// &pilot3,
		// &language1,
		// &language2,
		// &language3,
		// &pilotLanguage1,
		// &pilotLanguage2,
		// &pilotLanguage3,
		// &jet1,
		// &jet2,

	}
)
