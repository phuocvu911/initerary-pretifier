package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	isColor := false
	argCount := 0
	for _, a := range args {
		if a == "--color" { //special flag for formatting case
			isColor = true
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
	lookupFile := args[2]

	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Input not found")
		return
	}

	lookup, err := os.ReadFile(lookupFile)
	if err != nil {
		fmt.Println("Airport lookup not found")
	}

	//lookup malform

	//writefile
}
