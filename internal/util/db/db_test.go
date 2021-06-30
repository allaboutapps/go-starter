package db_test

import (
	"context"
	"database/sql"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/models"
	"allaboutapps.dev/aw/go-starter/internal/test"
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

			return newUser2.Insert(ctx, tx, boil.Infer())
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
			//nolint:errcheck
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
