package cmd

import (
	"bufio"
	"log"
	"os"

	"github.com/solloball/test_task/parser"
)

func Execute() {
	args := os.Args

	if len(args) != 2 {
		log.Fatal("uncorrect args length: should be 2")
	}

	file, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	parser.Parse(scanner)
}
