package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 1 && args[0] == "-h" {
		fmt.Println(usage)
		return
	}

	//isColor := false
	argCount := 0
	for _, a := range args {
		if a == "--color" {
			//isColor = true
		} else {
			argCount++
		}
	}
	if argCount != 3 {
		fmt.Println(usage)
		return
	}
}
