package core

type ActionPrototypes interface {
	After(*Message, *Flow) error
	Run(*Message, *Flow) error
	Before(*Message, *Flow) error
	Type() string
}

