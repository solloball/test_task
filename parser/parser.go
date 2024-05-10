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
	state, err := initState(scanner)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(state.Start)
	fmt.Println(state.Finish)
	fmt.Println(state.Price)
	for _, i := range state.IsEmpty {
		fmt.Println(i)
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
	for i := range res.IsEmpty {
		res.IsEmpty[i] = true
	}

	// parse date.
	if scanner.Scan() == false {
		return res, errors.New("incorrect file structure")
	}
	field := strings.Fields(scanner.Text())
	if len(field) != 2 {
		return res, errors.New(scanner.Text())
	}
	const layout = "15:04"
	res.Start, err = time.Parse(layout, field[0])
	if err != nil {
		return res, errors.New(scanner.Text())
	}

	res.Finish, err = time.Parse(layout, field[1])
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
