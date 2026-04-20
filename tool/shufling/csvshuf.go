package main

import (
	"encoding/csv"
	"log"
	"math/rand/v2" // use v2 for Shuffle, v1 is deprecated
	"os"
)

func main() {
	file, err := os.Open("airport-lookup.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	//find number of columns, and create a permutation of column indices
	numCols := len(records[0])
	p := make([]int, numCols)
	for i := range numCols {
		p[i] = i
	}
	rand.Shuffle(len(p), func(i, j int) {
		p[i], p[j] = p[j], p[i]
	})

	for i, row := range records {
		if len(row) < len(p) {
			log.Fatalf("row %d has too few columns: %d", i, len(row))
		}

		newRow := make([]string, len(p))
		for j, idx := range p {
			newRow[j] = row[idx]
		}
		records[i] = newRow
	}

	desfile, err := os.Create("airport-dynamic.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer desfile.Close()

	writer := csv.NewWriter(desfile)
	defer writer.Flush()

	if err := writer.WriteAll(records); err != nil {
		log.Fatal(err)
	}
}
