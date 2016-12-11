package core

//ActionPrototypes
type Task struct {}

func (m *Task) After(message *Message, flow *Flow) (err error) {
	//log.Println("Task.after: ", message)
	return
}

func (m *Task) Run(message *Message, flow *Flow) (err error) {
	//log.Println("Task.run: ", message)
	return
}

func (m *Task) Before(message *Message, flow *Flow) (err error) {
	//log.Println("Task.before: ", message)
	return
}

func (m *Task) Type() string {
	return  "Task"
}