package main

import (
	bt "github.com/derzombiiie/bytetrap"

	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		printALot()
	}
}

// does not return
func printALot() {
	err := bt.Write(os.Stdout)
	log.Fatalf("Error: %s", err)
}
