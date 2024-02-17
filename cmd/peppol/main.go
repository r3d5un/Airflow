package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r3d5un/Airflow/internal/data"
	"github.com/r3d5un/Airflow/internal/peppol"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("starting")

	dsn := flag.String(
		"dsn",
		"postgres://postgres:postgres@localhost:5432/warehouse?sslmode=disable",
		"Database DSN (postgres://user:password@localhost:5432/db)",
	)
	workers := flag.Int(
		"workers",
		25,
		"Sets DB connections and number of concurrent processes",
	)
	flag.Parse()

	logger.Info("opening db")
	db, err := openDB(*dsn, *workers)
	if err != nil {
		logger.Error("Error opening db", err)
		os.Exit(1)
	}
	defer db.Close()
	queryTimeout := time.Duration(10) * time.Second
	models := data.NewModels(db, logger, &queryTimeout)

	logger.Info("downloading peppol business cards")
	resp, err := http.Get("https://directory.peppol.eu/export/businesscards")
	if err != nil {
		logger.Error("Error downloading peppol business cards", err)
		os.Exit(1)
	}

	logger.Info("unmarshalling peppol business cards")
	var peppolData peppol.Root

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading response body", err)
		os.Exit(1)
	}

	err = xml.Unmarshal(d, &peppolData)
	if err != nil {
		logger.Error("Error unmarshalling xml", err)
		os.Exit(1)
	}

	logger.Info("ingesting peppol business cards")
	var wg sync.WaitGroup
	bcChan := make(chan data.PeppolBusinessCard, len(peppolData.BusinessCard))

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range bcChan {
				result, err := models.PeppolBusinessCards.Upsert(context.Background(), &u)
				if err != nil {
					logger.Error("error upserting peppol business card", err)
					return
				}

				logger.Info("upserted peppol business card", "id", result.ID)
			}
		}()
	}

	for _, bc := range peppolData.BusinessCard {
		bcChan <- data.PeppolBusinessCard{
			ID:                 bc.Participant.Value,
			Name:               bc.Entity.Name.Name,
			CountryCode:        bc.Entity.CountryCode,
			PeppolBusinessCard: &bc,
		}
	}
	close(bcChan)

	wg.Wait()

	logger.Info("all business cards upserted")
}

func openDB(dsn string, workers int) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(workers)
	db.SetMaxIdleConns(workers)

	duration, err := time.ParseDuration("5m")
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
