package models

import (
	"path"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"path/filepath"
	"mime/multipart"
	"time"
	"io"
	"os"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"../../lib/common"
)

type Image struct {
	Id           	int64  			`orm:"pk;auto;column(id)" json:"id"`
	Thumb        	string			`orm:"" json:"thumb"`
	Url        	string			`orm:"-" json:"url"`
	Image        	string			`orm:"" json:"image"`
	MimeType       	string			`orm:"" json:"mime_type"`
	Title       	string			`orm:"" json:"title"`
	Size       	int64			`orm:"" json:"size"`
	Name        	string			`orm:"" json:"name"`
	Created_at	time.Time		`orm:"auto_now_add;type(datetime);column(created_at)" json:"created_at"`

}

func (i *Image) TableName() string {
	return beego.AppConfig.String("db_images")
}

func init() {
	orm.RegisterModel(new(Image))
}

// AddImage insert a new Image into database and returns
// last inserted Id on success.
func AddImage(m *Image) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetImageById retrieves Image by Id. Returns error if
// Id doesn't exist
func GetImageById(id int64) (v *Image, err error) {
	o := orm.NewOrm()
	v = &Image{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllImage retrieves all Image matches certain condition. Returns empty list if
// no records exist
func GetAllImage(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Image))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Image
	qs = qs.OrderBy(sortFields...)
	objects_count, err := qs.Count()
	if err != nil {
		return
	}
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		meta = &map[string]int64{
			"objects_count": objects_count,
			"limit":         limit,
			"offset":        offset,
		}
		return ml, meta, nil
	}
	return nil, nil, err
}

// UpdateImage updates Image by Id and returns error if
// the record to be updated doesn't exist
func UpdateImageById(m *Image) (err error) {
	o := orm.NewOrm()
	v := Image{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteImage deletes Image by Id and returns error if
// the record to be deleted doesn't exist
func DeleteImage(id int64) (err error) {
	o := orm.NewOrm()
	v := Image{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Image{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func UploadImages(files map[string][]*multipart.FileHeader) (fileList []*Image, errs []error) {

	errs = []error{}
	for _, f := range files {

		//for each fileheader, get a handle to the actual file
		file, err := f[0].Open()
		defer file.Close()
		if err != nil {
			errs = append(errs, err)
			return
		}

		// rename & save
		name := common.Strtomd5(common.RandomString(10))
		ext := strings.ToLower(path.Ext(f[0].Filename))
		newname := fmt.Sprintf("%s%s", name, ext)

		//create destination file making sure the path is writeable.
		dir := common.GetFullPath(name)
		os.MkdirAll(dir, os.ModePerm)
		dst, err := os.Create(filepath.Join(dir, newname))
		defer dst.Close()
		if err != nil {
			errs = append(errs, err)
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			errs = append(errs, err)
			return
		}

		size, _ := common.GetFileSize(filepath.Join(dir, newname))
		newFile := &Image{
			Size: size,
			MimeType: f[0].Header.Get("Content-Type"),
			Image: newname,
			Name: f[0].Filename,
		}
		fileList = append(fileList, newFile)
	}

	_, errs = AddMultipleFiles(fileList)

	return
}

// AddMultipleFiles
// Use a prepared statement to increase inserting speed with multiple inserts.
func AddMultipleFiles(images []*Image) (ids []int64, errs []error) {

	o := orm.NewOrm()
	qs := o.QueryTable(&Image{})
	i, _ := qs.PrepareInsert()
	for _, image := range images {
		id, err := i.Insert(image)
		if err != nil {
			errs = append(errs, err)
		} else {
			ids = append(ids, id)
		}
	}
	// PREPARE INSERT INTO user (`name`, ...) VALUES (?, ...)
	// EXECUTE INSERT INTO user (`name`, ...) VALUES ("slene", ...)
	// EXECUTE ...
	// ...
	i.Close() // Don't forget to close the statement

	return
}

// GetImageFilterList
//
func GetImageFilterList() (ml []orm.Params, err error) {

	o := orm.NewOrm()
	image := &Image{}
	o.Raw(`
SELECT
	DATE_FORMAT(f.created_at, '%Y-%m-%d') as date, COUNT( f.created_at) as count
FROM `+ image.TableName()+` f
GROUP BY date(f.created_at)
ORDER BY f.created_at`).Values(&ml)

	return
}

// GetAllImagesByDate
//
func GetAllImagesByDate(filter string) (images []Image, err error) {

	o := orm.NewOrm()
	image := &Image{}
	o.Raw(`
SELECT *
FROM `+ image.TableName()+` f
WHERE DATE_FORMAT(f.created_at, '%Y %m %d') = DATE_FORMAT('`+filter+`', '%Y %m %d')
ORDER BY f.created_at`).QueryRows(&images)

	return
}

func (m *Image) GetUrl() {
	m.Url = common.GetLinkPath(m.Image)
}