package main

import (
	"log"

	"github.com/solloball/test_task/cmd"
)

func main() {
	setup()

	cmd.Execute()
}

func setup() {
	// Disable any other output in logger
	log.SetFlags(0)
}
