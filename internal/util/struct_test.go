package util_test

import (
	"context"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type insertable interface {
	Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
}

type testStructEmpty struct {
}

type testStructPrivateFiled struct {
	privateUser *models.User
}

type testStructPrimitives struct {
	X int
	Y string
}

type testStructFixture struct {
	User1               *models.User
	User2               *models.User
	User1AppUserProfile *models.AppUserProfile
	User1AccessToken1   *models.AccessToken

	X           int
	Y           string
	privateUser *models.User
}

func TestGetFieldsImplementingInvalidInput(t *testing.T) {

	// invalid interfaceObject input param, must be a pointer to an interface
	// pointer to a struct
	_, err := util.GetFieldsImplementing(&testStructEmpty{}, &testStructEmpty{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "interfaceObject")
	// pointer to a pointer to an interface
	interfaceObjPtr := (*insertable)(nil)
	_, err = util.GetFieldsImplementing(&testStructEmpty{}, &interfaceObjPtr)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "interfaceObject")

	// invalid structPtr input param, must be a pointer to a struct
	_, err = util.GetFieldsImplementing(testStructEmpty{}, (*insertable)(nil))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "structPtr")
	_, err = util.GetFieldsImplementing((*insertable)(nil), (*insertable)(nil))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "structPtr")
	_, err = util.GetFieldsImplementing([]*testStructEmpty{}, (*insertable)(nil))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "structPtr")
}

func TestGetFieldsImplementingNoFields(t *testing.T) {
	// No fields returned from empty structs
	structEmptyFields, err := util.GetFieldsImplementing(&testStructEmpty{}, (*insertable)(nil))
	assert.NoError(t, err)
	assert.Empty(t, structEmptyFields)

	// No fields returned from structs with only private fields
	structPrivate := testStructPrivateFiled{privateUser: &models.User{ID: "bfc9d3be-a13c-4790-befb-573c9a5b11a4"}}
	structPrivateFields, err := util.GetFieldsImplementing(&structPrivate, (*insertable)(nil))
	assert.NoError(t, err)
	assert.Empty(t, structPrivateFields)

	// No fields returned if struct fields are primitive
	structPrimitive := testStructPrimitives{X: 12, Y: "y"}
	structPrimitiveFields, err := util.GetFieldsImplementing(&structPrimitive, (*insertable)(nil))
	assert.NoError(t, err)
	assert.Empty(t, structPrimitiveFields)

	// No fieds returned if an interface is not matching
	type notMatchedInterface interface {
		// columns param missing
		Insert(ctx context.Context, exec boil.ContextExecutor) error
	}
	fix := testStructFixture{
		User1: &models.User{ID: "bfc9d3be-a13c-4790-befb-573c9a5b11a4"},
	}
	fixFields, err := util.GetFieldsImplementing(&fix, (*notMatchedInterface)(nil))
	assert.NoError(t, err)
	assert.Empty(t, fixFields)
}

func TestGetFieldsImplementingSuccess(t *testing.T) {
	// Struct not initialized
	// It's a responsibility of a user to make sure that the fields are not nil before using them.
	structNotInitialized := testStructFixture{}
	structNotInitializedFields, err := util.GetFieldsImplementing(&structNotInitialized, (*insertable)(nil))
	assert.NoError(t, err)
	assert.Equal(t, 4, len(structNotInitializedFields))
	for _, f := range structNotInitializedFields {
		assert.Nil(t, f)
		assert.Implements(t, (*insertable)(nil), f)
	}

	// Struct initialized
	fix := testStructFixture{
		privateUser:         &models.User{ID: "bfc9d3be-a13c-4790-befb-573c9a5b11a4"},
		User1:               &models.User{ID: "9e16c597-2491-45bb-89ca-775b6e07f51d"},
		User2:               &models.User{ID: "52028fd6-e299-4d36-8bba-21fe4713ffcd"},
		User1AppUserProfile: &models.AppUserProfile{UserID: "9e16c597-2491-45bb-89ca-775b6e07f51d"},
		User1AccessToken1:   &models.AccessToken{UserID: "9e16c597-2491-45bb-89ca-775b6e07f51d"},
		X:                   12,
		Y:                   "y",
	}

	insertableFields, err := util.GetFieldsImplementing(&fix, (*insertable)(nil))
	assert.NoError(t, err)
	assert.Equal(t, 4, len(insertableFields))

	for _, f := range insertableFields {
		assert.NotNil(t, f)
		assert.Implements(t, (*insertable)(nil), f)
	}

	type upsertable interface {
		Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error
	}
	upsertableFields, err := util.GetFieldsImplementing(&fix, (*upsertable)(nil))
	assert.NoError(t, err)
	// there should be equal number of fields implementing Insertable and Upsertable interface
	assert.Equal(t, len(insertableFields), len(upsertableFields))
	for _, f := range upsertableFields {
		assert.NotNil(t, f)
		assert.Implements(t, (*upsertable)(nil), f)
	}
}
