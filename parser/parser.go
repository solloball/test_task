package parser

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type State struct {
	Start   time.Time
	Finish  time.Time
	IsEmpty []bool
	Price   int64
}

func Parse(scanner *bufio.Scanner) {
	_, err := initState(scanner)
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func initState(scanner *bufio.Scanner) (State, error) {
	var res State

	// parse count tables.
	if scanner.Scan() == false {
		return res, errors.New("incorrect file structure")
	}
	n, err := strconv.ParseInt(scanner.Text(), 10, 32)
	if err != nil {
		return res, errors.New(scanner.Text())
	}
	res.IsEmpty = make([]bool, n)

	// parse date.
	if scanner.Scan() == false {
		return res, errors.New("incorrect file structure")
	}
	field := strings.Fields(scanner.Text())
	if len(field) != 2 {
		return res, errors.New(scanner.Text())
	}
	res.Start, err = time.Parse("00:00", field[0])
	if err != nil {
		return res, errors.New(scanner.Text())
	}

	res.Finish, err = time.Parse("00:00", field[1])
	if err != nil {
		return res, errors.New(scanner.Text())
	}

	// parse price.
	if scanner.Scan() == false {
		return res, errors.New(scanner.Text())
	}
	res.Price, err = strconv.ParseInt(scanner.Text(), 10, 32)
	if err != nil {
		return res, errors.New(scanner.Text())
	}

	return res, nil
}
