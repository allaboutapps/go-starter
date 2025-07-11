package db

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TxFn func(boil.ContextExecutor) error

func WithTransaction(ctx context.Context, db *sql.DB, txHandler TxFn) error {
	return WithConfiguredTransaction(ctx, db, nil, txHandler)
}

func WithConfiguredTransaction(ctx context.Context, db *sql.DB, options *sql.TxOptions, txHandler TxFn) error {
	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		util.LogFromContext(ctx).Warn().Err(err).Msg("Failed to start transaction")
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		cause := recover()

		switch {
		case cause != nil:
			util.LogFromContext(ctx).Error().Interface("cause", cause).Msg("Recovered from panic, rolling back transaction and panicking again")

			if txErr := tx.Rollback(); txErr != nil {
				util.LogFromContext(ctx).Warn().Err(txErr).Msg("Failed to roll back transaction after recovering from panic")
			}

			panic(cause)
		case err != nil:
			util.LogFromContext(ctx).Warn().Err(err).Msg("Received error, rolling back transaction")

			if txErr := tx.Rollback(); txErr != nil {
				util.LogFromContext(ctx).Warn().Err(txErr).Msg("Failed to roll back transaction after receiving error")
			}
		default:
			err = tx.Commit()
			if err != nil {
				util.LogFromContext(ctx).Warn().Err(err).Msg("Failed to commit transaction")
			}
		}
	}()

	err = txHandler(tx)
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	return nil
}

func NullIntFromInt64Ptr(i *int64) null.Int {
	if i == nil {
		return null.NewInt(0, false)
	}

	return null.NewInt(int(*i), true)
}

func NullFloat32FromFloat64Ptr(f *float64) null.Float32 {
	if f == nil {
		return null.NewFloat32(0.0, false)
	}

	return null.NewFloat32(float32(*f), true)
}

func NullIntFromInt16Ptr(i *int16) null.Int {
	if i == nil {
		return null.NewInt(0, false)
	}

	return null.NewInt(int(*i), true)
}

func Int16PtrFromNullInt(i null.Int) *int16 {
	if !i.Valid || i.Int > math.MaxInt16 || i.Int < math.MinInt16 {
		return nil
	}

	res := int16(i.Int)
	return &res
}

func Int16PtrFromInt(i int) *int16 {
	if i > math.MaxInt16 || i < math.MinInt16 {
		return nil
	}

	res := int16(i)
	return &res
}

func NullStringIfEmpty(s string) null.String {
	if len(s) == 0 {
		return null.String{}
	}

	return null.StringFrom(s)
}
