package core

const (
	STARTED		= Status("started")
	IN_PROCESS	= Status("in process")
	ENDED		= Status("ended")
	DONE		= Status("done")
)

type Status string
