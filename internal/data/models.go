package data

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Units UnitModel
}

func NewModels(db *sql.DB, logger *slog.Logger, timeout *time.Duration) Models {
	return Models{
		Units: UnitModel{DB: db, Logger: logger, Timeout: timeout},
	}
}
