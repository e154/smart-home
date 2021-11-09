// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package endpoint

import (
	"github.com/e154/smart-home/cmd/cli/commands/generate"
	"github.com/e154/smart-home/common"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
	"text/template"
)

var (
	log = common.MustGetLogger("endpoint")
)

var endpointTpl = `//CODE GENERATED AUTOMATICALLY

package {{.Package}}

import (
	"{{.Dir}}/common"
	m "{{.Dir}}/models"
	"{{.Dir}}/system/validation"
)

// {{.Name}} ...
type {{.Name}} struct {
	*CommonEndpoint
}

// New{{.Name}} ...
func New{{.Name}}(common *CommonEndpoint) *{{.Name}} {
	return &{{.Name}}{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *{{.Name}}) Add(params m.{{.ModelName}}) (result m.{{.ModelName}}, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.{{.AdaptorName}}.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.{{.AdaptorName}}.GetById(id)

	return
}

// GetByName ...
func (n *{{.Name}}) GetById(id int64) (result m.{{.ModelName}}, err error) {

	result, err = n.adaptors.{{.AdaptorName}}.GetById(id)

	return
}

// Update ...
func (n *{{.Name}}) Update(params m.{{.ModelName}}) (result m.{{.ModelName}}, errs []*validation.Error, err error) {

	var user m.{{.ModelName}}
	if user, err = n.adaptors.{{.AdaptorName}}.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&user, &params, common.JsonEngine)

	// validation
	_, errs = user.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.{{.AdaptorName}}.Update(user); err != nil {
		return
	}

	result, err = n.adaptors.{{.AdaptorName}}.GetById(user.Id)

	return
}

// GetList ...
func (n *{{.Name}}) GetList(limit, offset int64, order, sortBy string) (result []m.{{.ModelName}}, total int64, err error) {
	result, total, err = n.adaptors.{{.AdaptorName}}.List(limit, offset, order, sortBy)
	return
}

// Delete ...
func (n *{{.Name}}) Delete(id int64) (err error) {

	if _, err = n.adaptors.{{.AdaptorName}}.GetById(id); err != nil {
		return
	}

	err = n.adaptors.{{.AdaptorName}}.Delete(id)

	return
}

// Search ...
func (n *{{.Name}}) Search(query string, limit, offset int) (result []m.{{.ModelName}}, total int64, err error) {
	
	//result, total, err = n.adaptors.{{.AdaptorName}}.Search(query, limit, offset)
	
	return
}
`

var (
	endpointCmd = &cobra.Command{
		Use:   "e",
		Short: "endpoint generator",
		Long:  "$ cli g e [-m=ModelName, -a=AdaptorName] [endpointName]",
	}
	modelName   string
	adaptorName string
	packageName = "endpoint"
)

func init() {
	generate.Generate.AddCommand(endpointCmd)
	endpointCmd.Flags().StringVarP(&modelName, "model", "m", "interface{}", "interface{}")
	endpointCmd.Flags().StringVarP(&adaptorName, "adaptor", "a", "interface{}", "interface{}")
	endpointCmd.Run = func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			log.Error("Wrong number of arguments. Run: cli help generate")
			return
		}

		currpath, _ := os.Getwd()

		g := Generator{}
		g.Generate(args[0], currpath)
	}
}

type Generator struct{}

func (e Generator) Generate(endpointName, currpath string) {

	log.Infof("Using '%sEndpoint' as endpoint name", endpointName)
	log.Infof("Using '%s' as package name", packageName)

	fp := path.Join(currpath, "endpoint")

	e.addEndpoint(fp, endpointName)
}

func (e Generator) addEndpoint(fp, endpointName string) {

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		// Create the endpoint's directory
		if err := os.MkdirAll(fp, 0777); err != nil {
			log.Errorf("Could not create endpoints directory: %s", err.Error())
			return
		}
	}

	fpath := path.Join(fp, strings.ToLower(endpointName)+".go")
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Errorf("Could not create endpoint file: %s", err.Error())
		return
	}
	defer f.Close()

	templateData := struct {
		Package     string
		Name        string
		List        string
		ModelName   string
		AdaptorName string
		Dir         string
	}{
		Dir:         common.Dir(),
		Package:     packageName,
		Name:        endpointName,
		ModelName:   modelName,
		AdaptorName: adaptorName,
	}
	t := template.Must(template.New("endpoint").Parse(endpointTpl))

	if t.Execute(f, templateData) != nil {
		log.Error(err.Error())
	}

	common.FormatSourceCode(fpath)
}
