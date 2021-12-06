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
	log = common.MustGetLogger("db")
)

var dbModelTpl = `//CODE GENERATED AUTOMATICALLY

package {{.Package}}

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

// {{.Name}}s ...
type {{.Name}}s struct {
	Db *gorm.DB
}

// {{.Name}} ...
type {{.Name}} struct {
	Id                int64` + "`gorm:\"primary_key\"`" + `
	Name              string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// TableName ...
func (d *{{.Name}}) TableName() string {
	return "{{.Name}}"
}

// Add ...
func (n {{.Name}}s) Add(ver *{{.Name}}) (id int64, err error) {
	if err = n.Db.Create(&ver).Error; err != nil {
		return
	}
	id = ver.Id
	return
}

// GetById ...
func (n {{.Name}}s) GetById(id int64) (ver *{{.Name}}, err error) {
	ver = &{{.Name}}{Id: id}
	err = n.Db.First(&ver).Error
	return
}

// Update ...
func (n {{.Name}}s) Update(m *{{.Name}}) (err error) {
	q := map[string]interface{}{
		"name":        m.Name,
	}
	err = n.Db.Model(&{{.Name}}{Id: m.Id}).Updates(q).Error
	return
}

// Delete ...
func (n {{.Name}}s) Delete(id int64) (err error) {
	err = n.Db.Delete(&{{.Name}}{Id: id}).Error
	return
}

// List ...
func (n *{{.Name}}s) List(limit, offset int64, orderBy, sort string) (list []*{{.Name}}, total int64, err error) {

	if err = n.Db.Model({{.Name}}{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*{{.Name}}, 0)
	q := n.Db.Model(&{{.Name}}{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error

	return
}

// Search ...
func (n *{{.Name}}s) Search(query string, limit, offset int) (list []*{{.Name}}, total int64, err error) {

	q := n.Db.Model(&{{.Name}}{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*{{.Name}}, 0)
	err = q.Find(&list).Error

	return
}

`

var (
	dbModelCmd = &cobra.Command{
		Use:   "dbm",
		Short: "db model generator",
		Long:  "$ cli g dbm [Name]",
	}
	packageName = "db"
)

func init() {
	generate.Generate.AddCommand(dbModelCmd)
	dbModelCmd.Run = func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			log.Error("Wrong number of arguments. Run: cli help generate")
			return
		}

		currpath, _ := os.Getwd()

		g := Generator{}
		g.Generate(args[0], currpath)
	}
}

// Generator ...
type Generator struct{}

// Generate ...
func (e Generator) Generate(modelName, currpath string) {

	log.Infof("Using '%s' as model name", modelName)
	log.Infof("Using '%s' as package name", packageName)

	fp := path.Join(currpath, "db")

	e.addModel(fp, modelName)
}

func (e Generator) addModel(fp, modelName string) {

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		// Create the model's directory
		if err := os.MkdirAll(fp, 0777); err != nil {
			log.Errorf("Could not create db directory: %s", err.Error())
			return
		}
	}

	fpath := path.Join(fp, strings.ToLower(modelName)+".go")
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		log.Errorf("Could not create model file: %s", err.Error())
		return
	}
	defer f.Close()

	templateData := struct {
		Package     string
		Name        string
		List        string
		ModelName   string
		AdaptorName string
	}{
		Package: packageName,
		Name:    modelName,
	}
	t := template.Must(template.New("dbModel").Parse(dbModelTpl))

	if t.Execute(f, templateData) != nil {
		log.Error(err.Error())
	}

	common.FormatSourceCode(fpath)
}
