package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"

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
}
