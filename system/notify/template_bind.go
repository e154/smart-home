package notify

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

// Javascript Binding
//
// IC.Template()
//	 .render('name', {'key':'val'})
//
type TemplateBind struct {
	adaptor *adaptors.Adaptors
}

func (t *TemplateBind) Render(templateName string, params map[string]interface{}) *m.TemplateRender {
	render, err := t.adaptor.Template.Render(templateName, params)
	if err != nil {
		return nil
	}
	return render
}
