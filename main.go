package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	colorOutput := false
	argCount := 0
	for _, a := range args {
		if a == "--color" { //special flag for formatting case
			colorOutput = true
		} else {
			argCount++
		}
	}
	if argCount != 3 || args[0] == "-h" {
		fmt.Println(usage)
		return
	}

	inputPath := args[0]
	outputPath := args[1]
	lookupPath := args[2]

	// Check input exists
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Input not found")
		} else {
			fmt.Println(err)
		}
		return
	}

	// Check lookup file exist and Load airport lookup

	airports, err := LoadAirportLookup(lookupPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Process the itinerary
	result, colorResult := processItinerary(string(inputData), airports)

	// Write output, even though writefile can stop mid-op, our test case guard that to return when airport lookup malformed and return earlier. so using os.WriteFile here is fine
	//filemode: owner, group , others. 4: read, 2:write, 1:execute. Since it just a textfile, no need to execute permission.
	if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
		fmt.Println("Output path invalid or you dont have permission to write.")
		return
	}

	// Bonus: print colored output to stdout (only if --color flag is set)
	if colorOutput {
		fmt.Print(colorResult)
	}
}
