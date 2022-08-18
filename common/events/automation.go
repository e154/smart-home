package events

// EventCallTaskTrigger ...
type EventCallTaskTrigger struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// EventCallTaskAction ...
type EventCallTaskAction struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// EventEnableTask ...
type EventEnableTask struct {
	Id int64 `json:"id"`
}

// EventDisableTask ...
type EventDisableTask struct {
	Id int64 `json:"id"`
}

// EventAddedTask ...
type EventAddedTask struct {
	Id int64 `json:"id"`
}

// EventRemoveTask ...
type EventRemoveTask struct {
	Id int64 `json:"id"`
}

// EventUpdateTask ...
type EventUpdateTask struct {
	Id int64 `json:"id"`
}
