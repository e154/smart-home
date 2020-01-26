package core

//ActionPrototypes
type Task struct{}

func (m *Task) After(flow *Flow) (err error) {
	//log.Infof("Task.after: %v", message)
	return
}

func (m *Task) Run(flow *Flow) (err error) {
	//log.Infof("Task.run: %v", message)
	return
}

func (m *Task) Before(flow *Flow) (err error) {
	//log.Infof("Task.before: %v", message)
	return
}

func (m *Task) Type() string {
	return "Task"
}
