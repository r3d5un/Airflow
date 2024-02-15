package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/r3d5un/Airflow/internal/brreg/units"
)

type UnitSlice []units.Unit

func main() {
	data, err := os.ReadFile("./enheter_alle.json.gz")
	if err != nil {
		fmt.Println("Error reading file", err)
		os.Exit(1)
	}

	reader := bytes.NewReader(data)
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		fmt.Println("Error creating gzip reader", err)
		os.Exit(1)
	}

	unzippedData, err := io.ReadAll(gzReader)
	if err != nil {
		fmt.Println("Error reading unzipped data", err)
		os.Exit(1)
	}

	var units UnitSlice
	err = json.Unmarshal(unzippedData, &units)

	for _, unit := range units {
		fmt.Println(unit.Navn)
	}

	dsn := "postgres://postgres:postgres@localhost:5432/warehouse?sslmode=disable"
	db, err := openDB(dsn)
	if err != nil {
		fmt.Println("Error opening db", err)
		os.Exit(1)
	}
	defer db.Close()
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
