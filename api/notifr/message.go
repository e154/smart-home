package notifr

import "github.com/e154/smart-home/api/models"

type Message interface {
	save()	(string, *models.Message, error)
	send()	error
}