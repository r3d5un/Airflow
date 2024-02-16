package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r3d5un/Airflow/internal/brreg"
	"github.com/r3d5un/Airflow/internal/data"
)

type UnitSlice []brreg.Unit

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("reading file")
	unzippedData, err := read("./enheter_alle.json.gz")
	if err != nil {
		fmt.Println("Error reading file", err)
		os.Exit(1)
	}

	logger.Info("marshalling data")
	var units UnitSlice
	err = json.Unmarshal(unzippedData, &units)

	logger.Info("opening db")
	dsn := "postgres://postgres:postgres@localhost:5432/warehouse?sslmode=disable"
	db, err := openDB(dsn)
	if err != nil {
		fmt.Println("Error opening db", err)
		os.Exit(1)
	}
	defer db.Close()

	queryTimeout := time.Duration(10) * time.Second

	models := data.NewModels(db, logger, &queryTimeout)

	logger.Info("upserting units")
	var wg sync.WaitGroup
	unitChan := make(chan data.Unit, len(units))

	for i := 0; i < 25; i++ {
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

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
