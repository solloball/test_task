package runner

type inboxId int
type outboxId int

const (
	errAlreadyInside   = "YouShallNotPass"
	errNonWorkingHours = "NotOpenYet"
	errUnknownClient   = "ClientUnknown"
	errBusy            = "PlaceIsBusy"
	errFreeTables      = "ICanWaitNoLonger!"
)

const layout = "15:04"

const (
	eventVisit inboxId = iota + 1
	eventSitTable
	eventWait
	eventLeave
)

const (
	outEventLeave outboxId = iota + 11
	outEventSit
	outEventErr
)
