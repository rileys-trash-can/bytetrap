package main

import (
	bt "github.com/rileys-trash-can/bytetrap"

	_ "embed"
	"os"
	"strconv"
)

//go:embed help.txt
var helptext []byte

func main() {
	if len(os.Args) == 1 {
		err := bt.Write(os.Stdout)
		os.Stdout.Write([]byte(err.Error()))
		os.Exit(1)
	}

	if os.Args[1] == "-1" || os.Args[1] == "--one" {
		os.Stdout.Write([]byte(bt.GetPasta()))

		return
	}

	if os.Args[1] == "-n" && len(os.Args) == 3 {
		i, err := strconv.ParseUint(os.Args[2], 10, 64)
		if err != nil {
			os.Stdout.Write([]byte(err.Error()))
			os.Exit(1)
		}

		for i += 0; i > 0; i-- {
			os.Stdout.Write([]byte((bt.GetPasta())))
		}

		return
	}

	os.Stdout.Write(helptext)
	os.Exit(1)
}
