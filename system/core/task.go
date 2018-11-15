package core

//ActionPrototypes
type Task struct{}

func (m *Task) After(message *Message, flow *Flow) (err error) {
	//log.Info("Task.after: ", message)
	return
}

func (m *Task) Run(message *Message, flow *Flow) (err error) {
	//log.Info("Task.run: ", message)
	return
}

func (m *Task) Before(message *Message, flow *Flow) (err error) {
	//log.Info("Task.before: ", message)
	return
}

func (m *Task) Type() string {
	return "Task"
}
