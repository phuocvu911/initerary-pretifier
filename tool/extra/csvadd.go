package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strconv"
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

	records[0] = append(records[0], "type", "elevation_ft")
	for i := 1; i < len(records); i++ {
		records[i] = append(records[i], randomType(), strconv.Itoa(randomElevation()))
	} //strconv is way more efficient than fmt.Sprintf for int to string conversion

	desfile, err := os.Create("airport-extra.csv")
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

// randomType returns a random airport type
var types = []string{"large", "medium", "small", "heliport"}

func randomType() string {
	return types[rand.Intn(len(types))]
}

// randomElevation returns a random elevation in feet
func randomElevation() int {
	return rand.Intn(5000) // 0–4999 ft
}
