package models

type ActionPrototypes interface {
	After(*Message, *Flow) error
	Run(*Message, *Flow) error
	Before(*Message, *Flow) error
	Type() string
}

type Action interface {
	Compare(*Message, *Flow) error
}