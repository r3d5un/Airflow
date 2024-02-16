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

type SubUnit struct {
	ID           string         `json:"id"`
	ParentID     string         `json:"parent_id"`
	LastUpdate   *time.Time     `json:"last_updated"`
	BRREGSubUnit *brreg.SubUnit `json:"brreg_unit"`
}

type SubUnitModel struct {
	Timeout *time.Duration
	DB      *sql.DB
	Logger  *slog.Logger
}

func (m *SubUnitModel) Get(ctx context.Context, id string) (*SubUnit, error) {
	stmt := `SELECT id, parent_id, last_updated, brreg_subunit
FROM brreg_subunits
WHERE id = $1;`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	var su SubUnit
	var brregSubUnit []byte

	m.Logger.InfoContext(qCtx, "querying subunit", "query", stmt, "id", id)
	row := m.DB.QueryRowContext(ctx, stmt, id)

	err := row.Scan(&su.ID, &su.LastUpdate, brregSubUnit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			m.Logger.InfoContext(ctx, "no records found", "query", stmt, "id", id)
			return nil, ErrNoRecord
		} else {
			m.Logger.InfoContext(ctx, "unable to query subunit", "query", stmt, "id", id)
			return nil, err
		}
	}
	m.Logger.InfoContext(ctx, "data retrieved")

	err = json.Unmarshal(brregSubUnit, &su.BRREGSubUnit)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling brreg_subunit", "error", err)
		return nil, err
	}

	return &su, nil
}

func (m *SubUnitModel) Insert(ctx context.Context, su *SubUnit) (*SubUnit, error) {
	query := `INSERT INTO brreg_subunits (
id, parent_id, last_updated, brreg_unit
)
VALUES ($1, $2, $3, $4)
RETURNING id, parent_id, last_updated, brreg_subunit;`

	jsonUnit, err := json.Marshal(su.BRREGSubUnit)
	if err != nil {
		return nil, err
	}

	args := []any{
		su.ID,
		su.ParentID,
		su.LastUpdate,
		jsonUnit,
	}

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var brregSubUnit []byte

	m.Logger.InfoContext(rCtx, "inserting subunit", "query", query, "args", args)
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&su.ID, &su.ParentID, &su.LastUpdate, &brregSubUnit,
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
	m.Logger.InfoContext(ctx, "subunit inserted", "id", su.ID)

	err = json.Unmarshal(brregSubUnit, &su.BRREGSubUnit)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling brreg_subunit", "error", err)
		return nil, err
	}

	return su, nil
}

func (m *SubUnitModel) Upsert(ctx context.Context, su *SubUnit) (*SubUnit, error) {
	query := `INSERT INTO brreg_subunits (id, parent_id, last_updated, brreg_subunit)
VALUES ($1, $2, NOW(), $3)
ON CONFLICT (id) DO UPDATE
SET
    last_updated = EXCLUDED.last_updated,
    parent_id = EXCLUDED.parent_id,
    brreg_subunit = EXCLUDED.brreg_subunit
RETURNING id, parent_id, last_updated, brreg_subunit;`

	jsonUnit, err := json.Marshal(su.BRREGSubUnit)
	if err != nil {
		return nil, err
	}

	args := []any{
		su.ID,
		su.ParentID,
		jsonUnit,
	}

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var brregUnit []byte

	m.Logger.InfoContext(rCtx, "upserting subunit", "query", query, "args", args)
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&su.ID, &su.ParentID, &su.LastUpdate, &brregUnit,
	)
	if err != nil {
		m.Logger.ErrorContext(
			ctx,
			"unable to insert subunit",
			"query",
			query,
			"args",
			args,
			"error",
			err,
		)
		return su, err
	}
	m.Logger.InfoContext(ctx, "subunit upserted", "id", su.ID)

	err = json.Unmarshal(brregUnit, &su.BRREGSubUnit)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling brreg_subunit", "error", err)
		return su, err
	}

	return su, nil
}

func (m *SubUnitModel) Update(ctx context.Context, su *SubUnit) (*SubUnit, error) {
	query := `UPDATE brreg_subunits
SET id = $1, parent_id = $2, last_updated = NOW(), brreg_unit = $3
WHERE id = $1
RETURNING id, parent_id, last_updated, brreg_subunit;
    `

	jsonUnit, err := json.Marshal(su.BRREGSubUnit)
	if err != nil {
		return nil, err
	}

	args := []any{
		su.ID,
		su.ParentID,
		jsonUnit,
	}

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var brregUnit []byte

	m.Logger.InfoContext(rCtx, "updating subunit", "query", query, "args", args)
	err = m.DB.QueryRowContext(ctx, query, args...).Scan(
		&su.ID, &su.ParentID, &su.LastUpdate, &brregUnit,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			m.Logger.InfoContext(ctx, "no records found", "query", query, "args", args)
			return su, ErrNoRecord
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
	m.Logger.InfoContext(rCtx, "updated brreg subunit", "subunit", su.BRREGSubUnit)

	err = json.Unmarshal(brregUnit, &su.BRREGSubUnit)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling brreg_subunit", "error", err)
		return su, err
	}

	return su, nil
}

func (m *SubUnitModel) Delete(ctx context.Context, id int) (*SubUnit, error) {
	query := `DELETE FROM brreg_subunits WHERE id = $1
RETURNING id, parent_id, last_updated, brreg_subunit;`

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var su *SubUnit
	var brregSubUnit []byte

	m.Logger.InfoContext(rCtx, "deleting subunit", "query", query, "id", id)
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&su.ID, &su.ParentID, &su.LastUpdate, &brregSubUnit,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			m.Logger.InfoContext(ctx, "no records found", "query", query, "id", id)
			return su, ErrNoRecord
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
	m.Logger.InfoContext(ctx, "unit deleted", "subunit", su)

	err = json.Unmarshal(brregSubUnit, &su.BRREGSubUnit)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling brreg_subunit", "error", err)
		return su, err
	}

	return su, nil
}
