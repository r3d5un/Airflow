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
	stmt := `SELECT id, name, countrycode, last_updated, business_cards
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