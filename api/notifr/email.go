package notifr

import "github.com/e154/smart-home/api/models"

type Email struct {
	List		string	//"Alice <alice@example.com>, Bob <bob@example.com>, Eve <eve@example.com>"
	Subject		string
	Body		string
	Template	string
	Params		map[string]interface{}
}

func NewEmail() (email *Email) {
	email = &Email{
		Params: make(map[string]interface{}),
	}

	return
}

func (e *Email) make() (err error) {

	var r *models.EmailRender

	if e.Template != "" {

		if r, err = e.render(); err != nil {
			return err
		}

		e.Subject = r.Subject
		e.Body = r.Body
	}

	err = e.save()

	return
}

func (e *Email) render() (r *models.EmailRender, err error) {

	return
}

func (e *Email) save() (err error) {

	return
}

func (e *Email) send() (err error) {

	return
}