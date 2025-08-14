package db_test

import (
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	"allaboutapps.dev/aw/go-starter/internal/test/fixtures"
	swaggerTypes "allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderBy(t *testing.T) {
	test.WithTestDatabase(t, func(sqlDB *sql.DB) {
		ctx := t.Context()
		fix := fixtures.Fixtures()

		_, err := fix.UserRequiresConfirmation.Delete(ctx, sqlDB)
		require.NoError(t, err)

		noUsername := models.User{
			Scopes: types.StringArray{"cms"},
		}

		upperUsername := models.User{
			Username: null.StringFrom("USER3@example.com"),
			Scopes:   types.StringArray{"cms"},
		}

		err = noUsername.Insert(ctx, sqlDB, boil.Infer())
		require.NoError(t, err)

		err = upperUsername.Insert(ctx, sqlDB, boil.Infer())
		require.NoError(t, err)

		users, err := models.Users(db.OrderBy(swaggerTypes.OrderDirAsc, models.TableNames.Users, models.UserColumns.Username)).All(ctx, sqlDB)
		require.NoError(t, err)
		require.NotEmpty(t, users)
		assert.Equal(t, upperUsername.ID, users[0].ID)
		assert.Equal(t, upperUsername.Username, users[0].Username)

		users, err = models.Users(db.OrderByLower(swaggerTypes.OrderDirAsc, models.TableNames.Users, models.UserColumns.Username)).All(ctx, sqlDB)
		require.NoError(t, err)
		require.NotEmpty(t, users)
		assert.Equal(t, fix.User1.ID, users[0].ID)
		assert.Equal(t, fix.User1.Username, users[0].Username)

		users, err = models.Users(db.OrderByWithNulls(swaggerTypes.OrderDirAsc, db.OrderByNullsFirst, models.TableNames.Users, models.UserColumns.Username)).All(ctx, sqlDB)
		require.NoError(t, err)
		require.NotEmpty(t, users)
		assert.Equal(t, noUsername.ID, users[0].ID)
		assert.Equal(t, noUsername.Username, users[0].Username)

		users, err = models.Users(db.OrderByLowerWithNulls(swaggerTypes.OrderDirDesc, db.OrderByNullsLast, models.TableNames.Users, models.UserColumns.Username)).All(ctx, sqlDB)
		require.NoError(t, err)
		require.NotEmpty(t, users)
		assert.Equal(t, fix.UserDeactivated.ID, users[0].ID)
		assert.Equal(t, fix.UserDeactivated.Username, users[0].Username)
		assert.Equal(t, upperUsername.ID, users[1].ID)
		assert.Equal(t, upperUsername.Username, users[1].Username)
	})
}
