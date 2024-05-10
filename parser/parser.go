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

const (
	errId              = "13"
	errAlreadyInside   = "YouShallNotPass"
	errNonWorkingHours = "NotOpenYet"
)

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

type event struct {
	time  time.Time
	id    inboxId
	price int
	body  []string
}

type usage struct {
	time time.Time
	name visitorName
}

type table struct {
	IsEmpty     bool
	usageTime   time.Duration
	sittingTime time.Time
}

type visitorName string

type state struct {
	queue    []visitorName
	visitors []visitorName
	start    time.Time
	finish   time.Time
	tables   []table
	price    int64
}

func Parse(scanner *bufio.Scanner) {
	state, err := initState(scanner)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(state.start.Format(layout))
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
	fmt.Println(state.finish.Format(layout))
}

func handleEvent(ev event, state *state) error {
	switch ev.id {
	case id_1:
		handler1(ev, state)
		break
	case id_2:
		handler2(ev, state)
		break
	case id_3:
		handler3(ev, state)
		break
	case id_4:
		handler4(ev, state)
		break
	}
	return nil
}

func handler1(ev event, state *state) {
	visitor := ev.body[0]
	for _, name := range state.visitors {
		if name == visitorName(visitor) {
			fmt.Println(ev.time.Format(layout) + " " + errId + "" + errAlreadyInside)
		}
	}

	if ev.time.Before(state.start) || ev.time.After(state.finish) {
		fmt.Println(ev.time.Format(layout) + " " + errId + " " + errNonWorkingHours)
	}

	state.visitors = append(state.visitors, visitorName(visitor))
}

func handler2(ev event, state *state) {
}

func handler3(ev event, state *state) {
}

func handler4(ev event, state *state) {
}

func parseEvent(input string) (event, error) {
	var res event
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

func initState(scanner *bufio.Scanner) (state, error) {
	var res state

	// parse count tables.
	if scanner.Scan() == false {
		return res, errors.New("incorrect file structure")
	}
	n, err := strconv.ParseInt(scanner.Text(), 10, 32)
	if err != nil {
		return res, errors.New(scanner.Text())
	}
	res.tables = make([]table, n)
	for i := range res.tables {
		res.tables[i] = table{true, time.Duration(0), time.Time{}}
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

	res.queue = make([]visitorName, 0)
	res.visitors = make([]visitorName, 0)

	return res, nil
}
