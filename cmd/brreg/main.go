package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r3d5un/Airflow/internal/brreg"
	"github.com/r3d5un/Airflow/internal/data"
)

type UnitSlice []brreg.Unit
type SubUnitSlice []brreg.SubUnit

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dataset := flag.String(
		"dataset",
		"units",
		"The desired dataset. Choice between units and subunits",
	)
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

	switch *dataset {
	case "units":
		units, err := downloadUnits(logger)
		if err != nil {
			logger.Error("Error downloading subunits", err)
			os.Exit(1)
		}
		ingestUnits(units, *workers, &models, logger)
	case "subunits":
		subunits, err := downloadSubUnits(logger)
		if err != nil {
			logger.Error("Error downloading subunits", err)
			os.Exit(1)
		}
		ingestSubUnits(subunits, *workers, &models, logger)
	default:
		logger.Error("no valid dataset selected")
		os.Exit(1)
	}

	logger.Info("done")
}

func downloadSubUnits(logger *slog.Logger) (SubUnitSlice, error) {
	logger.Info("downloading subunits")
	resp, err := http.Get("https://data.brreg.no/enhetsregisteret/oppslag/underenheter/lastned/v2")
	if err != nil {
		fmt.Println("Error downloading file", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	logger.Info("unzipping subunits")
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(d)
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	d, err = io.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}

	logger.Info("unmarshalling subunits")
	var subunits SubUnitSlice
	err = json.Unmarshal(d, &subunits)
	if err != nil {
		return nil, err
	}

	return subunits, nil
}

func downloadUnits(logger *slog.Logger) (UnitSlice, error) {
	logger.Info("downloading units")
	resp, err := http.Get("https://data.brreg.no/enhetsregisteret/oppslag/enheter/lastned/v2")
	if err != nil {
		fmt.Println("Error downloading file", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	logger.Info("unzipping subunits")
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(d)
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	d, err = io.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}

	logger.Info("unmarshalling subunits")
	var units UnitSlice
	err = json.Unmarshal(d, &units)
	if err != nil {
		return nil, err
	}

	return units, nil
}

func read(path string) (d []byte, err error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(file)
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	d, err = io.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func ingestUnits(units UnitSlice, workers int, models *data.Models, logger *slog.Logger) {
	logger.Info("upserting units")
	var wg sync.WaitGroup
	unitChan := make(chan data.Unit, len(units))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range unitChan {
				result, err := models.Units.Upsert(context.Background(), &u)
				if err != nil {
					logger.Error("error upserting unit", "error", err)
					return
				}

				logger.Info("unit upserted", "id", result.ID)
			}
		}()
	}

	for _, u := range units {
		record := data.Unit{
			ID:        u.OrganisasjonsNummer,
			BRREGUnit: &u,
		}
		unitChan <- record
	}
	close(unitChan)

	wg.Wait()

	logger.Info("all units upserted")
}

func ingestSubUnits(subunits SubUnitSlice, workers int, models *data.Models, logger *slog.Logger) {
	logger.Info("upserting subunits")
	var wg sync.WaitGroup
	subUnitChan := make(chan data.SubUnit, len(subunits))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for su := range subUnitChan {
				result, err := models.SubUnits.Upsert(context.Background(), &su)
				if err != nil {
					logger.Error("error upserting subunit", "error", err)
					return
				}

				logger.Info("subunit upserted", "id", result.ID)
			}
		}()
	}

	for _, su := range subunits {
		record := data.SubUnit{
			ID:           su.OrganisasjonsNummer,
			ParentID:     su.OverordnetEnhet,
			BRREGSubUnit: &su,
		}
		subUnitChan <- record
	}
	close(subUnitChan)

	wg.Wait()

	logger.Info("all subunits upserted")
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
