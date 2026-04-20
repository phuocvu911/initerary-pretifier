package main

import (
	"encoding/csv"
	"log"
	"os"
)

var p = []int{2, 4, 1, 5, 3, 0}

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

	out, err := os.Create("airport-dynamic.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	defer writer.Flush()

	if err := writer.WriteAll(records); err != nil {
		log.Fatal(err)
	}
}
