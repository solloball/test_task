package runner

import "time"

type event struct {
	time  time.Time
	id    inboxId
	price int
	body  []string
}
