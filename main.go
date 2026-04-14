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
}
