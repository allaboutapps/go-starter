package db_test

import (
	"context"
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
	swaggerTypes "allaboutapps.dev/aw/go-starter/internal/types"
	"allaboutapps.dev/aw/go-starter/internal/util/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"
)

func TestWithTransactionSuccess(t *testing.T) {
	test.WithTestDatabase(t, func(sqlDB *sql.DB) {
		ctx := context.Background()

		count, err := models.Users().Count(ctx, sqlDB)
		require.NoError(t, err)
		assert.Greater(t, count, int64(0))

		err = db.WithTransaction(ctx, sqlDB, func(tx boil.ContextExecutor) error {
			newUser := models.User{
				IsActive: true,
				Username: null.StringFrom("test"),
				Scopes:   types.StringArray{"cms"},
			}

			if err := newUser.Insert(ctx, tx, boil.Infer()); err != nil {
				return err
			}

			newCount, err := models.Users().Count(ctx, tx)
			require.NoError(t, err)
			assert.Equal(t, count+1, newCount)

			delCnt, err := models.Users().DeleteAll(ctx, tx)
			if err != nil {
				return err
			}
			assert.Equal(t, newCount, delCnt)

			newCount, err = models.Users().Count(ctx, tx)
			require.NoError(t, err)
			assert.Equal(t, int64(0), newCount)

			return nil
		})
		require.NoError(t, err)

		newCount, err := models.Users().Count(ctx, sqlDB)
		require.NoError(t, err)
		assert.Equal(t, int64(0), newCount)
	})
}

func TestWithTransactionWithError(t *testing.T) {
	test.WithTestDatabase(t, func(sqlDB *sql.DB) {
		ctx := context.Background()

		count, err := models.Users().Count(ctx, sqlDB)
		require.NoError(t, err)
		assert.Greater(t, count, int64(0))

		err = db.WithTransaction(ctx, sqlDB, func(tx boil.ContextExecutor) error {
			newUser := models.User{
				IsActive: true,
				Username: null.StringFrom("test"),
				Scopes:   types.StringArray{"cms"},
			}

			err := newUser.Insert(ctx, tx, boil.Infer())
			require.NoError(t, err)

			newCount, err := models.Users().Count(ctx, tx)
			require.NoError(t, err)
			assert.Equal(t, count+1, newCount)

			delCnt, err := models.Users().DeleteAll(ctx, tx)
			require.NoError(t, err)
			assert.Equal(t, newCount, delCnt)

			newCount, err = models.Users().Count(ctx, tx)
			require.NoError(t, err)
			assert.Equal(t, int64(0), newCount)

			newUser2 := models.User{
				IsActive: true,
				Username: null.StringFrom("test"),
			}
			if err := newUser2.Insert(ctx, tx, boil.Infer()); err != nil {
				return err
			}

			return nil
		})
		require.Error(t, err)

		newCount, err := models.Users().Count(ctx, sqlDB)
		require.NoError(t, err)
		assert.Equal(t, count, newCount)
	})
}

func TestWithTransactionWithPanic(t *testing.T) {
	test.WithTestDatabase(t, func(sqlDB *sql.DB) {
		ctx := context.Background()

		count, err := models.Users().Count(ctx, sqlDB)
		require.NoError(t, err)
		assert.Greater(t, count, int64(0))

		panicFunc := func() {
			_ = db.WithTransaction(ctx, sqlDB, func(tx boil.ContextExecutor) error {
				newUser := models.User{
					IsActive: true,
					Username: null.StringFrom("test"),
					Scopes:   types.StringArray{"cms"},
				}

				err := newUser.Insert(ctx, tx, boil.Infer())
				require.NoError(t, err)

				newCount, err := models.Users().Count(ctx, tx)
				require.NoError(t, err)
				assert.Equal(t, count+1, newCount)

				delCnt, err := models.Users().DeleteAll(ctx, tx)
				require.NoError(t, err)
				assert.Equal(t, newCount, delCnt)

				newCount, err = models.Users().Count(ctx, tx)
				require.NoError(t, err)
				assert.Equal(t, int64(0), newCount)

				panic("some panic")
			})
		}

		require.Panics(t, panicFunc)

		newCount, err := models.Users().Count(ctx, sqlDB)
		require.NoError(t, err)
		assert.Equal(t, count, newCount)
	})
}

func TestDBTypeConversions(t *testing.T) {
	i := int64(19)
	res := db.NullIntFromInt64Ptr(&i)
	assert.Equal(t, 19, res.Int)
	assert.True(t, res.Valid)

	res = db.NullIntFromInt64Ptr(nil)
	assert.False(t, res.Valid)

	f := 19.9999
	res2 := db.NullFloat32FromFloat64Ptr(&f)
	assert.Equal(t, float32(19.9999), res2.Float32)
	assert.True(t, res2.Valid)

	res2 = db.NullFloat32FromFloat64Ptr(nil)
	assert.False(t, res2.Valid)
}

func TestSearchStringToTSQuery(t *testing.T) {
	expected := "abcde:* & 12345:* & xyz:*"
	in := "    abcde 12345 xyz   "
	out := db.SearchStringToTSQuery(in)
	assert.Equal(t, expected, out)

	expected = "abcde:*"
	in = "abcde"
	out = db.SearchStringToTSQuery(in)
	assert.Equal(t, expected, out)
}

func TestInnterJoinWithFilter(t *testing.T) {
	test.WithTestDatabase(t, func(sqlDB *sql.DB) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		profiles, err := models.AppUserProfiles(db.InnerJoinWithFilter(models.TableNames.AppUserProfiles,
			models.AppUserProfileColumns.UserID,
			models.TableNames.Users,
			models.UserColumns.ID,
			models.UserColumns.Username,
			"user1@example.com",
		)).All(ctx, sqlDB)
		require.NoError(t, err)
		require.Len(t, profiles, 1)

		assert.Equal(t, fixtures.User1AppUserProfile.UserID, profiles[0].UserID)
	})
}

func TestInnterJoin(t *testing.T) {
	test.WithTestDatabase(t, func(sqlDB *sql.DB) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		profiles, err := models.AppUserProfiles(db.InnerJoin(models.TableNames.AppUserProfiles,
			models.AppUserProfileColumns.UserID,
			models.TableNames.Users,
			models.UserColumns.ID,
		),
			models.UserWhere.Username.EQ(null.StringFrom("user1@example.com")),
		).All(ctx, sqlDB)
		require.NoError(t, err)
		require.Len(t, profiles, 1)

		assert.Equal(t, fixtures.User1AppUserProfile.UserID, profiles[0].UserID)
	})
}

func TestOrderBy(t *testing.T) {
	test.WithTestDatabase(t, func(sqlDB *sql.DB) {
		ctx := context.Background()
		fixtures := test.Fixtures()

		noUsername := models.User{
			Scopes: types.StringArray{"cms"},
		}

		upperUsername := models.User{
			Username: null.StringFrom("USER3@example.com"),
			Scopes:   types.StringArray{"cms"},
		}

		err := noUsername.Insert(ctx, sqlDB, boil.Infer())
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
		assert.Equal(t, fixtures.User1.ID, users[0].ID)
		assert.Equal(t, fixtures.User1.Username, users[0].Username)

		users, err = models.Users(db.OrderByWithNulls(swaggerTypes.OrderDirAsc, db.OrderByNullsFirst, models.TableNames.Users, models.UserColumns.Username)).All(ctx, sqlDB)
		require.NoError(t, err)
		require.NotEmpty(t, users)
		assert.Equal(t, noUsername.ID, users[0].ID)
		assert.Equal(t, noUsername.Username, users[0].Username)

		users, err = models.Users(db.OrderByLowerWithNulls(swaggerTypes.OrderDirDesc, db.OrderByNullsLast, models.TableNames.Users, models.UserColumns.Username)).All(ctx, sqlDB)
		require.NoError(t, err)
		require.NotEmpty(t, users)
		assert.Equal(t, fixtures.UserDeactivated.ID, users[0].ID)
		assert.Equal(t, fixtures.UserDeactivated.Username, users[0].Username)
		assert.Equal(t, upperUsername.ID, users[1].ID)
		assert.Equal(t, upperUsername.Username, users[1].Username)
	})
}
