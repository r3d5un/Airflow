package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/r3d5un/Airflow/internal/brreg/units"
)

type Unit struct {
	ID         int        `json:"id"`
	LastUpdate *time.Time `json:"last_update"`
	BRREGUnit  units.Unit `json:"brreg_unit"`
}

type UnitModel struct {
	Timeout *time.Duration
	DB      *sql.DB
	Logger  *slog.Logger
}

func (m *UnitModel) Get(ctx context.Context, id int) (*Unit, error) {
	if id < 1 {
		m.Logger.InfoContext(ctx, "invalid id", "id", id)
		return nil, ErrRecordNotFound
	}
	stmt := `
SELECT id, last_update, brreg_unit
FROM brreg_units
WHERE id = $1;`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	m.Logger.InfoContext(qCtx, "querying unit", "query", stmt, "id", id)
	row := m.DB.QueryRowContext(ctx, stmt, id)
	u := &Unit{}

	err := row.Scan(
		&u.ID,
		&u.LastUpdate,
		&u.BRREGUnit,
	)
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

	return u, nil
}

func (m *UnitModel) Insert(ctx context.Context, u *Unit) (*Unit, error) {
	query := `INSERT INTO brreg_units (
id, last_update, brreg_unit
)
VALUES ($1, $2, $3)
RETURNING id, last_update, brreg_unit;`

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

	m.Logger.InfoContext(rCtx, "inserting unit", "query", query, "args", args)
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&u.ID, &u.LastUpdate, &u.BRREGUnit,
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

	return u, nil
}

func (m *UnitModel) Upsert(ctx context.Context, u *Unit) (*Unit, error) {
	query := `INSERT INTO posts (id, last_update, brreg_unit)
VALUES ($1, NOW(), $2)
ON CONFLICT (id) DO UPDATE
SET last_update = EXCLUDED.last_update, brreg_unit = EXCLUDED.brreg_unit
RETURNING id, last_update, brreg_unit;`

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

	m.Logger.InfoContext(rCtx, "inserting unit", "query", query, "args", args)
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&u.ID, &u.LastUpdate, &u.BRREGUnit,
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

	return u, nil
}

func (m *UnitModel) Update(ctx context.Context, u *Unit) (*Unit, error) {
	query := `UPDATE posts
SET id = $2, last_update = NOW(), brreg_unit = $4
WHERE id = $1
RETURNING id, last_update, brreg_unit;
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

	m.Logger.InfoContext(rCtx, "updating unit", "query", query, "args", args)
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&u.ID, &u.LastUpdate, &u.BRREGUnit,
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

	return u, nil
}

func (m *UnitModel) Delete(ctx context.Context, id int) (*Unit, error) {
	query := "DELETE FROM posts WHERE id = $1 RETURNING id, last_update, brreg_unit;"

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var u *Unit

	m.Logger.InfoContext(rCtx, "deleting unit", "query", query, "id", id)
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.LastUpdate, &u.BRREGUnit,
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

	return u, nil
}
