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

type inboxId int
type outboxId int

const layout = "15:04"

const (
	id_1 inboxId = iota
	id_2
	id_3
	id_4
)

const (
	id_11 outboxId = iota
	id_12
	id_13
)

type Event struct {
	time  time.Time
	id    inboxId
	price int
	body  []string
}

type Table struct {
	IsEmpty   bool
	usageTime time.Duration
}

type State struct {
	start  time.Time
	finish time.Time
	tables []Table
	price  int64
}

func Parse(scanner *bufio.Scanner) {
	state, err := initState(scanner)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(state.start) // make format
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		event, err := parseEvent(scanner.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := handleEvent(event, &state); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(state.finish) // make format
}

func handleEvent(event Event, state *State) error {
	return nil
}

func parseEvent(input string) (Event, error) {
	var res Event
	var err error
	fields := strings.Fields(input)
	if len(fields) != 3 && len(fields) != 4 {
		return res, errors.New(input)
	}

	res.time, err = time.Parse(layout, fields[0])
	if err != nil {
		return res, errors.New(input)
	}

	val, err := strconv.ParseInt(fields[1], 10, 32)
	if err != nil {
		return res, errors.New(input)
	}
	res.id = inboxId(val)

	res.body = fields[2:]
	return res, nil
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
	res.tables = make([]Table, n)
	for i := range res.tables {
		res.tables[i] = Table{true, time.Duration(0)}
	}

	// parse date.
	if scanner.Scan() == false {
		return res, errors.New("incorrect file structure")
	}
	field := strings.Fields(scanner.Text())
	if len(field) != 2 {
		return res, errors.New(scanner.Text())
	}
	res.start, err = time.Parse(layout, field[0])
	if err != nil {
		return res, errors.New(scanner.Text())
	}

	res.finish, err = time.Parse(layout, field[1])
	if err != nil {
		return res, errors.New(scanner.Text())
	}

	// parse price.
	if scanner.Scan() == false {
		return res, errors.New(scanner.Text())
	}
	res.price, err = strconv.ParseInt(scanner.Text(), 10, 32)
	if err != nil {
		return res, errors.New(scanner.Text())
	}

	return res, nil
}
