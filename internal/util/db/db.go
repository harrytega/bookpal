package db

import (
	"context"
	"database/sql"
	"math"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"test-project/internal/util"
)

type TxFn func(boil.ContextExecutor) error

func WithTransaction(ctx context.Context, db *sql.DB, fn TxFn) error {
	return WithConfiguredTransaction(ctx, db, nil, fn)
}

func WithConfiguredTransaction(ctx context.Context, db *sql.DB, options *sql.TxOptions, fn TxFn) error {
	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		util.LogFromContext(ctx).Warn().Err(err).Msg("Failed to start transaction")
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			util.LogFromContext(ctx).Error().Interface("p", p).Msg("Recovered from panic, rolling back transaction and panicking again")

			if txErr := tx.Rollback(); txErr != nil {
				util.LogFromContext(ctx).Warn().Err(txErr).Msg("Failed to roll back transaction after recovering from panic")
			}

			panic(p)
		} else if err != nil {
			util.LogFromContext(ctx).Warn().Err(err).Msg("Received error, rolling back transaction")

			if txErr := tx.Rollback(); txErr != nil {
				util.LogFromContext(ctx).Warn().Err(txErr).Msg("Failed to roll back transaction after receiving error")
			}
		} else {
			err = tx.Commit()
			if err != nil {
				util.LogFromContext(ctx).Warn().Err(err).Msg("Failed to commit transaction")
			}
		}
	}()

	err = fn(tx)

	return err
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
