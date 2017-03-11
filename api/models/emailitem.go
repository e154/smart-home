package models

import (
	"time"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Tree struct {
	Name			string		`json:"name"`
	Description		string		`json:"description"`
	Status			string		`json:"status"`
	Nodes			[]*Tree		`json:"nodes"`
}

type EmailItem struct {
	Name			string		`orm:"pk;size(64);column(name);unique" valid:"Required;MaxSize(64)" json:"name"`
	Description		string		`orm:"size(255)" json:"description"`
	Content			string		`orm:"" json:"content"`
	Status			string		`orm:"size(64)" valid:"Required;MaxSize(64)" json:"status"`	//active, inactive
	Type			string		`orm:"size(64)" valid:"Required;MaxSize(64)" json:"type"`		//item, template
	Parent			string		`orm:"size(64)" valid:"MaxSize(64)" json:"parent"`
	Markers			[]string	`orm:"-" json:"markers"`
	Created_at		time.Time	`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`
	Updated_at		time.Time	`orm:"auto_now;type(datetime);column(updated_at)" json:"updated_at"`
}

type Items []*EmailItem

func (i *EmailItem) TableName() string {
	return beego.AppConfig.String("db_email_templates")
}

func (i *EmailItem) GetTemplate() (tpl *EmailTemplate, err error) {

	tpl = new(EmailTemplate)
	err = json.Unmarshal([]byte(i.Content), tpl)
	return
}

func (i Items) Len() int {
	return len(i)
}

func (i Items) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i Items) Less(a, b int) bool {
	_, items, _ := EmailItemGetList()

	result_a := Items{}
	result_b := Items{}
	getItemParents(items, &result_a, i[a].Name)
	getItemParents(items, &result_b, i[b].Name)

	return len(result_a) < len(result_b)
}

func init() {
	orm.RegisterModel(
	new(EmailItem),
	)
}

func EmailItemAddNew(body []byte) (id int64, b bool, valid validation.Validation, err error) {

	item := new(EmailItem)
	if err = json.Unmarshal(body, item); err != nil {
		return
	}

	item.Type = "item"

	b, err = valid.Valid(item)
	if err != nil || !b {
		return
	}

	o := orm.NewOrm()
	id, err = o.Insert(item)

	return
}

func EmailItemGet(item_name string) (item *EmailItem, err error) {

	o := orm.NewOrm()
	item = new(EmailItem)
	err = o.QueryTable(item).Filter("name", item_name).One(item)
	return
}

func EmailItemUpdate(body []byte, item_name string) (id int64, b bool, valid validation.Validation, err error) {

	item, err := EmailItemGet(item_name)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, item); err != nil {
		return
	}

	b, err = valid.Valid(item)
	if err != nil || !b {
		return
	}

	params := orm.Params{
		"content": item.Content,
		"description": item.Description,
		"parent": item.Parent,
		"status": item.Status,
		"type": item.Type,
	}

	o := orm.NewOrm()
	id, err = o.QueryTable(item).Filter("name", item_name).Update(params)

	return
}

func EmailItemDelete(item_name string) (err error) {

	o := orm.NewOrm()

	//check if item has chilren
	item := new(EmailItem)
	count, err := o.QueryTable(item).Filter("parent", item_name).Count()
	if err != nil {
		return
	}

	if count > 0 {
		err = errors.New("item has children")
		return
	}

	_, err = o.QueryTable(item).Filter("name", item_name).Filter("type", "item").Delete()

	return
}

func EmailItemGetList() (count int64, items Items, err error) {

	o := orm.NewOrm()
	table := new(EmailItem)
	count, err = o.QueryTable(table).Filter("type", "item").All(&items)
	return
}

func EmailItemGetSortedList() (count int, new_items []string, err error) {

	o := orm.NewOrm()
	table := new(EmailItem)
	items := Items{}
	_, err = o.QueryTable(table).Filter("type", "item").Filter("status", "active").All(&items)

	treeGetEndPoints := func(i Items, t *[]string) {
		for _, v := range i {
			var exist bool
			for _, k := range i {
				if k.Parent == v.Name {
					exist = true
					break
				}
			}

			if !exist {
				*t = append(*t, v.Name)
			}
		}
	};

	treeGetEndPoints(items, &new_items)
	count = len(new_items)
	return
}

func renderTreeRecursive(i Items, t *Tree, c string) {

	for _, item := range i {
		if item.Parent == c {
			tree := new(Tree)
			tree.Name = item.Name
			tree.Description = item.Description
			tree.Nodes = make([]*Tree, 0)	// fix - nodes: null
			tree.Status = item.Status
			t.Nodes = append(t.Nodes, tree)
			renderTreeRecursive(i, tree, item.Name)
		}
	}

	return
}

func EmailItemGetTree() (tree *Tree, err error) {

	o := orm.NewOrm()
	items := Items{}
	table := new(EmailItem)
	if _, err = o.QueryTable(table).Filter("type", "item").All(&items); err != nil {
		return
	}

	tree = new(Tree)
	for _, item := range items {
		if item.Parent == "" {
			tree.Description = item.Description
			tree.Name = item.Name
		}
	}

	renderTreeRecursive(items, tree, tree.Name)

	return
}

func EmailItemParentUpdate(name, parent string) (err error) {

	o := orm.NewOrm()
	table := new(EmailItem)
	_, err = o.QueryTable(table).Filter("name", name).Update(orm.Params{"parent": parent,})

	return
}

func updateTreeRecursive(t []*Tree, parent string) {

	for _, v := range t {
		if parent != "" {
			go EmailItemParentUpdate(v.Name, parent)
		}
		updateTreeRecursive(v.Nodes, v.Name)
	}

}

func EmailItemUpdateTree(body []byte) (err error) {

	tree := []*Tree{}
	err = json.Unmarshal(body, &tree)
	if err != nil {
		return
	}

	updateTreeRecursive(tree, "")

	return
}

func getItemParents(items Items, result *Items, s string){

	for _, item := range items {
		if item.Name == s {
			var exist bool
			for _, v := range *result {
				if v.Name == item.Name {
					exist = true
				}
			}
			if !exist {
				*result = append(*result, item)
			}
			getItemParents(items, result, item.Parent)
		}
	}
}
