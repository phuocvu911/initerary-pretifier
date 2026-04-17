package main

import (
	"encoding/csv"
	"errors"
	"os"
	"strings"

	m "pretifier/model"
)

var malform = errors.New("Airport lookup malformed")

func LoadAirportLookup(path string) (map[string]m.AirportRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("Airport lookup not found")
		} else {
			return nil, err
		}
	}
	defer file.Close()
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll() //rows [][]string
	if err != nil {               //err in reading csv, not relate to malform criteria
		return nil, err
	}
	if len(rows) == 0 {
		return nil, malform //empty file -> missing cols -> malformed data
	}

	header := rows[0]
	//dynamic column order.
	colIndex := make(map[string]int)
	requiredCols := []string{"name", "iso_country", "municipality", "icao_code", "iata_code", "coordinates"}
	for i, h := range header {
		colIndex[strings.TrimSpace(strings.ToLower(h))] = i
	}

	for _, col := range requiredCols {
		if _, ok := colIndex[col]; !ok {
			return nil, malform //missing required column, exit early
		}
	}

	airports := make(map[string]m.AirportRecord)
	for _, row := range rows[1:] {
		//check blank cell
		for _, col := range requiredCols {
			idx := colIndex[col]
			if idx >= len(row) || strings.TrimSpace(row[idx]) == "" {
				return nil, malform
			}
		}

		rec := m.AirportRecord{
			Name:         strings.TrimSpace(row[colIndex["name"]]),
			Municipality: strings.TrimSpace(row[colIndex["municipality"]]),
			ICAOCode:     strings.TrimSpace(row[colIndex["icao_code"]]),
			IATACode:     strings.TrimSpace(row[colIndex["iata_code"]]),
		}

		airports[rec.ICAOCode] = rec
		airports[rec.IATACode] = rec
	}
	return airports, nil
}

func ParseISO8601(s string)