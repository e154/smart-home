package notifr

type Message interface {
	send()	error
}