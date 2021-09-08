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

package adaptor

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
	log = common.MustGetLogger("adaptor")
)

var adaptorTpl = `//CODE GENERATED AUTOMATICALLY

package {{.Package}}

import (
	"{{.Dir}}/common"
	"{{.Dir}}/db"
	m "{{.Dir}}/models"
	"github.com/jinzhu/gorm"
)

type I{{.Name}} interface {
	Add(ver *m.{{.ModelName}}) (id int64, err error)
	GetBy{{.Name}}Name(imageName string) (ver *m.{{.ModelName}}, err error)
	GetById(mapId int64) (ver *m.{{.ModelName}}, err error)
	Update(ver *m.{{.ModelName}}) (err error)
	Delete(mapId int64) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.{{.ModelName}}, total int64, err error)
	fromDb(*db.{{.ModelName}}) *m.{{.ModelName}}
	toDb(*m.{{.ModelName}}) *db.{{.ModelName}}
}

// {{.Name}} ...
type {{.Name}} struct {
	I{{.Name}}
	table *db.{{.Name}}s
	db    *gorm.DB
}

// Get{{.Name}}Adaptor ...
func Get{{.Name}}Adaptor(d *gorm.DB) I{{.Name}} {
	return &{{.Name}}{
		table: &db.{{.ModelName}}s{Db: d},
		db:    d,
	}
}

// Add ...
func (n *{{.Name}}) Add(ver *m.{{.ModelName}}) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

// GetById ...
func (n *{{.Name}}) GetById(mapId int64) (ver *m.{{.ModelName}}, err error) {

	var dbVer *db.{{.ModelName}}
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *{{.Name}}) Update(ver *m.{{.ModelName}}) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

// Delete ...
func (n *{{.Name}}) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

// List ...
func (n *{{.Name}}) List(limit, offset int64, orderBy, sort string) (list []*m.{{.ModelName}}, total int64, err error) {

	if sort == "" {
		sort = "id"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.{{.ModelName}}
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.{{.ModelName}}, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *{{.Name}}) fromDb(dbVer *db.{{.ModelName}}) (ver *m.{{.ModelName}}) {
	ver = &m.{{.ModelName}}{
		Id:        dbVer.Id,
		UpdatedAt: dbVer.UpdatedAt,
		CreatedAt: dbVer.CreatedAt,
	}
	return
}

func (n *{{.Name}}) toDb(ver *m.{{.ModelName}}) (dbVer *db.{{.ModelName}}) {
	dbVer = &db.{{.ModelName}}{
		Id:       ver.Id,
		CreatedAt:     ver.CreatedAt,
		UpdatedAt:     ver.UpdatedAt,
	}
	return
}

`

var (
	adaptorCmd = &cobra.Command{
		Use:   "a",
		Short: "adaptor generator",
		Long:  "$ cli g a [-m=ModelName] [adaptorName]",
	}
	modelName   string
	adaptorName string
	packageName = "adaptors"
)

func init() {
	generate.Generate.AddCommand(adaptorCmd)
	adaptorCmd.Flags().StringVarP(&modelName, "model", "m", "interface{}", "interface{}")
	adaptorCmd.Run = func(cmd *cobra.Command, args []string) {

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

func (e Generator) Generate(adaptorName, currpath string) {

	log.Infof("Using '%s' as adaptor name", adaptorName)
	log.Infof("Using '%s' as package name", packageName)

	fp := path.Join(currpath, "adaptors")

	e.addAdaptor(fp, adaptorName)
}

func (e Generator) addAdaptor(fp, adaptorName string) {

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		// Create the adaptor's directory
		if err := os.MkdirAll(fp, 0777); err != nil {
			log.Errorf("Could not create adaptors directory: %s", err.Error())
			return
		}
	}

	fpath := path.Join(fp, strings.ToLower(adaptorName)+".go")
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Errorf("Could not create adaptor file: %s", err.Error())
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
		Name:        adaptorName,
		ModelName:   modelName,
		AdaptorName: adaptorName,
	}
	t := template.Must(template.New("adaptor").Parse(adaptorTpl))

	if t.Execute(f, templateData) != nil {
		log.Error(err.Error())
	}

	common.FormatSourceCode(fpath)
}
