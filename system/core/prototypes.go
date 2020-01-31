package core

type ActionPrototypes interface {
	After(*Flow) error
	Run(*Flow) error
	Before(*Flow) error
	Type() string
}

