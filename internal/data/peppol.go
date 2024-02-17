package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/r3d5un/Airflow/internal/peppol"
)

type PeppolBusinessCard struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	CountryCode        string               `json:"countrycode"`
	LastUpdated        *time.Time           `json:"last_updated"`
	PeppolBusinessCard *peppol.BusinessCard `json:"peppol_business_card"`
}

type PeppolBusinessCardModel struct {
	Timeout *time.Duration
	DB      *sql.DB
	Logger  *slog.Logger
}

func (m *PeppolBusinessCardModel) Get(ctx context.Context, id string) (*PeppolBusinessCard, error) {
	stmt := `SELECT id, name, countrycode, last_updated, business_card
FROM peppol_business_cards
WHERE id = $1;`

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	var bc PeppolBusinessCard
	var jsonBc []byte

	m.Logger.InfoContext(qCtx, "querying peppol business card", "query", stmt, "id", id)
	row := m.DB.QueryRowContext(ctx, stmt, id)
	err := row.Scan(&bc.ID, &bc.Name, &bc.CountryCode, &bc.LastUpdated, jsonBc)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			m.Logger.InfoContext(ctx, "no records found", "query", stmt, "id", id)
			return nil, ErrNoRecord
		} else {
			m.Logger.InfoContext(ctx, "unable to query business card", "query", stmt, "id", id)
			return nil, err
		}
	}
	m.Logger.InfoContext(ctx, "data retrieved")

	err = json.Unmarshal(jsonBc, &bc.PeppolBusinessCard)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling peppol_business_card", "error", err)
		return nil, err
	}

	return &bc, nil

}

func (m *PeppolBusinessCardModel) Insert(ctx context.Context, bc *PeppolBusinessCard) (*PeppolBusinessCard, error) {
	stmt := `INSERT INTO peppol_business_cards (id, name, countrycode, last_updated, business_card)
VALUES ($1, $2, $3, NOW(), $5)
RETURNING id, name, countrycode, last_updated, business_card;`

	bcRaw, err := json.Marshal(bc.PeppolBusinessCard)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error marshaling peppol_business_card", "error", err)
		return nil, err
	}

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	var jsonBc []byte

	m.Logger.InfoContext(qCtx, "inserting peppol business card", "query", stmt, "id", bc.ID)
	row := m.DB.QueryRowContext(ctx, stmt, bc.ID, bc.Name, bc.CountryCode, bcRaw)
	err = row.Scan(&bc.ID, &bc.Name, &bc.CountryCode, &bc.LastUpdated, &jsonBc)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error inserting peppol business card", "error", err)
		return nil, err
	}
	m.Logger.InfoContext(ctx, "peppol business card inserted", "id", bc.ID)

	err = json.Unmarshal(jsonBc, &bc.PeppolBusinessCard)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling peppol_business_card", "error", err)
		return nil, err
	}

	return bc, nil
}

func (m *PeppolBusinessCardModel) Upsert(ctx context.Context, bc *PeppolBusinessCard) (*PeppolBusinessCard, error) {
	stmt := `INSERT INTO peppol_business_cards (id, name, countrycode, last_updated, business_card)
VALUES ($1, $2, $3, NOW(), $5)
ON CONFLICT (id) DO UPDATE
SET
    name = EXCLUDED.name, countrycode = EXCLUDED.countrycode,
    last_updated = NOW(), business_card = EXCLUDED.business_card
RETURNING id, name, countrycode, last_updated, business_card;`

	bcRaw, err := json.Marshal(bc.PeppolBusinessCard)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error marshaling peppol_business_card", "error", err)
		return nil, err
	}

	qCtx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	var jsonBc []byte

	m.Logger.InfoContext(qCtx, "upserting peppol business card", "query", stmt, "id", bc.ID)
	row := m.DB.QueryRowContext(ctx, stmt, bc.ID, bc.Name, bc.CountryCode, bcRaw)
	err = row.Scan(&bc.ID, &bc.Name, &bc.CountryCode, &bc.LastUpdated, &jsonBc)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error upserting peppol business card", "error", err)
		return nil, err
	}
	m.Logger.InfoContext(ctx, "peppol business card upserted", "id", bc.ID)

	err = json.Unmarshal(jsonBc, &bc.PeppolBusinessCard)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling peppol_business_card", "error", err)
		return nil, err
	}

	return bc, nil
}

func (m *PeppolBusinessCardModel) Update(ctx context.Context, bc *PeppolBusinessCard) (*PeppolBusinessCard, error) {
	stmt := `UPDATE peppol_business_card
SET id = $1, name = $2, countrycode = $3, last_updated = NOW(), business_card = $4
WHERE id = $1
RETURNING id, name, countrycode, last_updated, business_card;
    `

	jsonBc, err := json.Marshal(bc.PeppolBusinessCard)
	if err != nil {
		return nil, err
	}

	args := []any{
		bc.ID,
		bc.Name,
		bc.CountryCode,
		jsonBc,
	}

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var bcRaw []byte

	m.Logger.InfoContext(rCtx, "updating business card", "query", stmt, "args", args)
	err = m.DB.QueryRowContext(ctx, stmt, args...).Scan(
		&bc.ID, &bc.Name, &bc.CountryCode, &bc.LastUpdated, &bcRaw,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			m.Logger.InfoContext(ctx, "no records found", "query", stmt, "args", args)
			return nil, ErrNoRecord
		default:
			m.Logger.ErrorContext(
				ctx, "unable to update business card",
				"query", stmt, "args", args, "error", err,
			)
			return nil, err
		}
	}
	m.Logger.InfoContext(rCtx, "updated brreg business card", "business_card", bc.PeppolBusinessCard)

	err = json.Unmarshal(bcRaw, &bc.PeppolBusinessCard)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling business card", "error", err)
		return nil, err
	}

	return bc, nil
}

func (m *PeppolBusinessCardModel) Delete(ctx context.Context, id string) (*PeppolBusinessCard, error) {
	query := `DELETE FROM peppol_business_card
WHERE id = $1
RETURNING id, name, countrycode, last_updated, business_card;`

	rCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var bc *PeppolBusinessCard
	var bcRaw []byte

	m.Logger.InfoContext(rCtx, "deleting business card", "query", query, "id", id)
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&bc.ID, &bc.LastUpdated, &bcRaw,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			m.Logger.InfoContext(ctx, "no records found", "query", query, "id", id)
			return bc, ErrNoRecord
		default:
			m.Logger.ErrorContext(
				ctx, "unable to update business card",
				"query", query, "id", id, "error", err,
			)
			return nil, err
		}
	}
	m.Logger.InfoContext(ctx, "business card deleted", "business_card", bc)

	err = json.Unmarshal(bcRaw, &bc.PeppolBusinessCard)
	if err != nil {
		m.Logger.ErrorContext(ctx, "error unmarshaling business card", "error", err)
		return bc, err
	}

	return bc, nil
}
