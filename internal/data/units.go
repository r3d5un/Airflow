package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/r3d5un/Airflow/internal/brreg"
)

type Unit struct {
	ID         string      `json:"id"`
	LastUpdate *time.Time  `json:"last_updated"`
	BRREGUnit  *brreg.Unit `json:"brreg_unit"`
}

type UnitModel struct {
	Timeout *time.Duration
	DB      *sql.DB
	Logger  *slog.Logger
}

func (m *UnitModel) Get(ctx context.Context, id string) (*Unit, error) {
	stmt := `SELECT id, last_updated, brreg_unit
FROM brreg_units
WHERE id = $1;`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	var u Unit
	var brregUnit []byte

	m.Logger.InfoContext(qCtx, "querying unit", "query", stmt, "id", id)
	row := m.DB.QueryRowContext(ctx, stmt, id)

	err := row.Scan(&u.ID, &u.LastUpdate, brregUnit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			m.Logger.InfoContext(ctx, "no records found", "query", stmt, "id", id)
			return nil, ErrNoRecord
		} else {
			m.Logger.InfoContext(ctx, "unable to query unit", "query", stmt, "id", id)
			return nil, err
		}
	}
	m.Logger.InfoContext(ctx, "data retrieved")

	err = json.Unmarshal(brregUnit, &u.BRREGUnit)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling brreg_unit", "error", err)
		return nil, err
	}

	return &u, nil
}

func (m *UnitModel) Insert(ctx context.Context, u *Unit) (*Unit, error) {
	query := `INSERT INTO brreg_units (
id, last_updated, brreg_unit
)
VALUES ($1, $2, $3)
RETURNING id, last_updated, brreg_unit;`

	jsonUnit, err := json.Marshal(u.BRREGUnit)
	if err != nil {
		return nil, err
	}

	args := []any{
		u.ID,
		u.LastUpdate,
		jsonUnit,
	}

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var brregUnit []byte

	m.Logger.InfoContext(rCtx, "inserting unit", "query", query, "args", args)
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&u.ID, &u.LastUpdate, &brregUnit,
	)
	if err != nil {
		m.Logger.ErrorContext(
			ctx,
			"unable to insert unit",
			"query",
			query,
			"args",
			args,
			"error",
			err,
		)
		return nil, err
	}
	m.Logger.InfoContext(ctx, "unit upserted", "id", u.ID)

	err = json.Unmarshal(brregUnit, &u.BRREGUnit)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling brreg_unit", "error", err)
		return nil, err
	}

	return u, nil
}

func (m *UnitModel) Upsert(ctx context.Context, u *Unit) (*Unit, error) {
	query := `INSERT INTO brreg_units (id, last_updated, brreg_unit)
VALUES ($1, NOW(), $2)
ON CONFLICT (id) DO UPDATE
SET last_updated = EXCLUDED.last_updated, brreg_unit = EXCLUDED.brreg_unit
RETURNING id, last_updated, brreg_unit;`

	jsonUnit, err := json.Marshal(u.BRREGUnit)
	if err != nil {
		return nil, err
	}

	args := []any{
		u.ID,
		jsonUnit,
	}

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var brregUnit []byte

	m.Logger.InfoContext(rCtx, "upserting unit", "query", query, "args", args)
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&u.ID, &u.LastUpdate, &brregUnit,
	)
	if err != nil {
		m.Logger.ErrorContext(
			ctx,
			"unable to insert unit",
			"query",
			query,
			"args",
			args,
			"error",
			err,
		)
		return u, err
	}
	m.Logger.InfoContext(ctx, "unit inserted", "id", u.ID)

	err = json.Unmarshal(brregUnit, &u.BRREGUnit)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling brreg_unit", "error", err)
		return u, err
	}

	return u, nil
}

func (m *UnitModel) Update(ctx context.Context, u *Unit) (*Unit, error) {
	query := `UPDATE brreg_units
SET id = $1, last_updated = NOW(), brreg_unit = $2
WHERE id = $1
RETURNING id, last_updated, brreg_unit;
    `

	jsonUnit, err := json.Marshal(u.BRREGUnit)
	if err != nil {
		return nil, err
	}

	args := []any{
		u.ID,
		jsonUnit,
	}

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var brregUnit []byte

	m.Logger.InfoContext(rCtx, "updating unit", "query", query, "args", args)
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&u.ID, &u.LastUpdate, &brregUnit,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			m.Logger.InfoContext(ctx, "no records found", "query", query, "args", args)
			return u, ErrNoRecord
		default:
			m.Logger.ErrorContext(
				ctx,
				"unable to update unit",
				"query",
				query,
				"args",
				args,
				"error",
				err,
			)
			return nil, err
		}
	}
	m.Logger.InfoContext(rCtx, "updated brreg unit", "unit", u.BRREGUnit)

	err = json.Unmarshal(brregUnit, &u.BRREGUnit)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling brreg_unit", "error", err)
		return u, err
	}

	return u, nil
}

func (m *UnitModel) Delete(ctx context.Context, id int) (*Unit, error) {
	query := "DELETE FROM brreg_units WHERE id = $1 RETURNING id, last_updated, brreg_unit;"

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var u *Unit
	var brregUnit []byte

	m.Logger.InfoContext(rCtx, "deleting unit", "query", query, "id", id)
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.LastUpdate, &brregUnit,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			m.Logger.InfoContext(ctx, "no records found", "query", query, "id", id)
			return u, ErrNoRecord
		default:
			m.Logger.ErrorContext(
				ctx,
				"unable to update unit",
				"query",
				query,
				"id",
				id,
				"error",
				err,
			)
			return nil, err
		}
	}
	m.Logger.InfoContext(ctx, "unit deleted", "unit", u)

	err = json.Unmarshal(brregUnit, &u.BRREGUnit)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling brreg_unit", "error", err)
		return u, err
	}

	return u, nil
}
