package data

import (
	"context"
	"database/sql"
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
