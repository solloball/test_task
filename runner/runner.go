package runner

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Run(scanner *bufio.Scanner) {
	state, err := parseInitState(scanner)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(state.start.Format(layout))

	for scanner.Scan() {
		fmt.Println(scanner.Text())

		event, err := parseEvent(scanner.Text())
		if err != nil {
			log.Println(err)
			continue
		}

		handleEvent(event, &state)
	}

	removeVisitors(&state)
	fmt.Println(state.finish.Format(layout))
	finish(state)
}

func finish(state state) {
	for i, t := range state.tables {
		log.Println(
			fmt.Sprint(i+1),
			fmt.Sprint(t.paidHours*state.price),
			time.Time{}.Add(t.usageTime).Format(layout),
		)
	}
}

func removeVisitors(state *state) {
	// make sorted slice
	sorted := make([]string, len(state.visitors))
	for i, k := range state.visitors {
		sorted[i] = string(k)
	}
	sort.Strings(sorted)

	for _, vis := range sorted {
		for i, t := range state.tables {
			if !t.IsEmpty && t.visitor == visitorName(vis) {
				state.tables[i].leave(state.finish)
				break
			}
		}
		printOutComeEvent(outEventLeave, state.finish, vis)
	}
}

func handleEvent(ev event, state *state) {
	switch ev.id {
	case eventVisit:
		handleVisit(ev, state)
		break
	case eventSitTable:
		handleSit(ev, state)
		break
	case eventWait:
		handleWait(ev, state)
		break
	case eventLeave:
		handleLeave(ev, state)
		break
	}
}

func handleVisit(ev event, state *state) {
	visitor := ev.body[0]
	for _, name := range state.visitors {
		if name == visitorName(visitor) {
			printOutComeEvent(outEventErr, ev.time, errAlreadyInside)
			return
		}
	}

	if ev.time.Before(state.start) || ev.time.After(state.finish) {
		printOutComeEvent(outEventErr, ev.time, errNonWorkingHours)
		return
	}

	state.visitors = append(state.visitors, visitorName(visitor))
}

func handleSit(ev event, state *state) {
	name := ev.body[0]
	table, _ := strconv.ParseInt(ev.body[1], 10, 32)

	if !find(state.visitors, visitorName(name)) {
		printOutComeEvent(outEventErr, ev.time, errUnknownClient)
		return
	}

	if !state.tables[table-1].IsEmpty {
		printOutComeEvent(outEventErr, ev.time, errBusy)
		return
	}

	state.tables[table-1].sit(ev.time, visitorName(name))
}

func find(slice []visitorName, name visitorName) bool {
	for _, visitor := range slice {
		if visitor == name {
			return true
		}
	}
	return false
}

func handleWait(ev event, state *state) {
	for _, t := range state.tables {
		if t.IsEmpty {
			printOutComeEvent(outEventErr, ev.time, errFreeTables)
			return
		}
	}

	name := ev.body[0]
	if len(state.queue) >= len(state.tables) {
		printOutComeEvent(outEventLeave, ev.time, ev.body[0])
	}

	state.queue = append(state.queue, visitorName(name))
}

func handleLeave(ev event, state *state) {
	name := visitorName(ev.body[0])
	if !find(state.visitors, visitorName(name)) {
		printOutComeEvent(outEventErr, ev.time, errUnknownClient)
		return
	}

	for i, vis := range state.visitors {
		if vis == name {
			state.visitors = remove(state.visitors, i)
			break
		}
	}

	for i, vis := range state.queue {
		if vis == name {
			state.queue = remove(state.queue, i)
			break
		}
	}

	for i, t := range state.tables {
		if !t.IsEmpty && t.visitor == name {
			state.tables[i].leave(ev.time)

			if len(state.queue) > 0 {
				state.tables[i].sit(ev.time, state.queue[0])
				printOutComeEvent(
					outEventSit,
					ev.time,
					string(state.queue[0])+" "+fmt.Sprint(i+1))
				state.queue = state.queue[1:]
			}
		}
	}
}

func remove(slice []visitorName, s int) []visitorName {
	return append(slice[:s], slice[s+1:]...)
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

func parseInitState(scanner *bufio.Scanner) (state, error) {
	var res state

	// parse count tables.
	if scanner.Scan() == false {
		return res, errors.New("empty file")
	}
	n, err := strconv.ParseInt(scanner.Text(), 10, 32)
	if err != nil {
		return res, errors.New(scanner.Text())
	}
	res.tables = make([]table, n)
	for i := range res.tables {
		res.tables[i] = table{
			true,
			time.Duration(0),
			time.Time{},
			0,
			visitorName("")}
	}

	// parse date.
	if scanner.Scan() == false {
		return res, errors.New(scanner.Text())
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
	pr, err := strconv.ParseInt(scanner.Text(), 10, 32)
	res.price = int(pr)
	if err != nil {
		return res, errors.New(scanner.Text())
	}

	res.queue = make([]visitorName, 0)
	res.visitors = make([]visitorName, 0)

	return res, nil
}

func printOutComeEvent(id outboxId, t time.Time, body string) {
	log.Println(t.Format(layout), int(id), body)
}
