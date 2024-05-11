package runner

import (
	"time"
)

type usage struct {
	time time.Time
	name visitorName
}

type table struct {
	IsEmpty     bool
	usageTime   time.Duration
	sittingTime time.Time
	paidHours   int
	visitor     visitorName
}

type visitorName string

type state struct {
	queue    []visitorName
	visitors []visitorName
	start    time.Time
	finish   time.Time
	tables   []table
	price    int
}

// Table should be empty in time calling
func (t *table) sit(tm time.Time, name visitorName) {
	t.sittingTime = tm
	t.visitor = name
	t.IsEmpty = false
}

// Table should be bussy in time calling
func (t *table) leave(tm time.Time) {
	duration := tm.Sub(t.sittingTime)

	t.paidHours += int((duration.Hours())) + 1
	t.usageTime += duration
	t.IsEmpty = true
}
