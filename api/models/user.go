package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"time"
	"unicode/utf8"
	"github.com/e154/smart-home/api/common"
	"encoding/json"
)

type User struct {
	Id                     	int64                  	`orm:"pk;auto" json:"id"`
	Nickname               	string                	`orm:"size(255)" valid:"Required;MinSize(3);MaxSize(255)" json:"nickname"`
	FirstName             	string                	`orm:"size(255)" valid:"MaxSize(255)" json:"first_name"`
	LastName              	string                	`orm:"size(255)" valid:"MaxSize(255)" json:"last_name"`
	EncryptedPassword	string                	`orm:"size(255)" valid:"Required;MaxSize(255)" json:"-"`
	Email                  	string                	`orm:"size(255)" valid:"Required;Email" json:"email"`
	HistoryStr              string                	`orm:"column(history)" json:"-"`
	History                	[]*UserHistory          `orm:"-" json:"history"`
	Status                 	string                	`orm:"size(255);default(blocked)" valid:"MaxSize(255)" json:"status"` //active, blocked
	ResetPasswordToken   	string                	`orm:"size(255)" json:"-"`
	AuthenticationToken   	string                	`orm:"size(255)" valid:"MaxSize(255)" json:"-"`
	Avatar                 	*Image                	`orm:"rel(fk);null;column(image_id)" json:"avatar"`
	SignInCount          	int64                   `orm:"size(11)" json:"sign_in_count"`
	CurrentSignInIp     	string                	`orm:"size(255);default(null)" json:"current_sign_in_ip"`
	LastSignInIp        	string                	`orm:"size(255);default(null)" json:"last_sign_in_ip"`

	CreatedBy	       	*User                   `orm:"rel(fk);null;column(user_id)" json:"created_by"`
	Role                   	*Role                	`orm:"rel(fk);null;column(role_name)" json:"role"`
	Meta			[]*UserMeta		`orm:"reverse(many);null" json:"meta"`

	ResetPasswordSentAt 	time.Time        	`orm:"type(datetime)" json:"-"`
	CurrentSignInAt     	time.Time        	`orm:"type(datetime);null;default(null)" json:"current_sign_in_at"`
	LastSignInAt        	time.Time        	`orm:"type(datetime);null;default(null)" json:"last_sign_in_at"`
	Created_at            	time.Time        	`orm:"auto_now_add;type(datetime)" json:"created_at"`
	Update_at            	time.Time       	`orm:"auto_now;type(datetime)" json:"update_at"`
	Deleted           	time.Time        	`orm:"type(datetime);null;default(null)" json:"deleted"`
}

const HISTORY_MAX = 8

func (m *User) TableName() string {
	return beego.AppConfig.String("db_users")
}

func init() {
	orm.RegisterModel(new(User))
}

// AddUser insert a new User into database and returns
// last inserted Id on success.
func AddUser(m *User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUserById(id int64) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Id: id}
	if err = o.Read(v); err != nil {
		return
	}

	err = v.GetHistory()

	return
}

// GetAllUser retrieves all User matches certain condition. Returns empty list if
// no records exist
func GetAllUser(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, meta *map[string]int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
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

	var l []User
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

// UpdateUser updates User by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserById(m *User) (err error) {
	o := orm.NewOrm()
	v := User{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUser(id int64) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&User{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func (u *User) LoadRelated() (err error) {
	o := orm.NewOrm()

	_, err = o.LoadRelated(u, "CreatedBy")
	_, err = o.LoadRelated(u, "Role")
	_, err = o.LoadRelated(u, "Meta")

	if u.CreatedBy != nil {
		_, err = o.LoadRelated(u, "CreatedBy")
	}

	if u.Avatar != nil {
		_, err = o.LoadRelated(u, "Avatar")
		u.Avatar.GetUrl()
	}

	return
}

// SetMeta
//
func (u *User) SetMeta(key, value string) (err error) {
	err = SetUserMeta(u, key, value)
	return
}

// GetMeta
//
func (u *User) GetMeta(key string) string {
	return GetUserMeta(u, key)
}

//
// etc
//
func (u *User) UpdateHistory(t time.Time, ipv4 string) {

	history := []*UserHistory{}
	json.Unmarshal([]byte(u.HistoryStr), &history)
	l := len(history)
	// max HISTORY_MAX records
	if l > HISTORY_MAX {
		// удаление элемента из начала среза
		history = history[l - HISTORY_MAX:]
	}

	history = append(history, &UserHistory{Ip: ipv4, Time: t})
	b, _ := json.Marshal(history)
	u.HistoryStr = string(b)
}

func (u *User) GetHistory() (err error) {

	err = json.Unmarshal([]byte(u.HistoryStr), &u.History)

	return
}

func (u *User) SignIn(ipv4 string) {

	// update count
	u.SignInCount += 1

	// update time
	last_t := u.CurrentSignInAt
	current_t := time.Now()

	u.LastSignInAt = last_t
	u.CurrentSignInAt = current_t

	// update ipv4
	last_ip := u.CurrentSignInIp
	current_ip := ipv4
	u.LastSignInIp = last_ip
	u.CurrentSignInIp = current_ip

	u.UpdateHistory(current_t, current_ip)

	o := orm.NewOrm()
	o.Update(u)
}

func (u *User) NewToken() (token string, err error) {

	num, err := beego.AppConfig.Int("auth_token_size")
	if err != nil {
		num = 50
	}

	o := orm.NewOrm()
	qs := o.QueryTable(u)

	// check token unique
	for {
		token = common.RandStr(num, "alphanum")
		u.AuthenticationToken = token

		if exist := qs.Filter("authentication_token", token).Exist(); !exist {
			break
		}
	}

	_, err = o.Update(u, "authentication_token")

	return
}

func (u *User) SetToken(token string) error {

	o := orm.NewOrm()
	table := new(User)
	qs := o.QueryTable(table)
	if len(token) != 0 {
		if exist := qs.Filter("authentication_token", token).Exist(); exist {
			return errors.New("Token not unique");
		}
	}

	u.AuthenticationToken = token
	_ , err := o.Update(u, "authentication_token")

	return err
}

func UserGetByAuthenticationToken(token string) (user *User, err error) {

	if utf8.RuneCountInString(token) > 255 {
		return
	}

	o := orm.NewOrm()
	user = new(User)
	user.AuthenticationToken = token
	if err = o.Read(user, "authentication_token"); err != nil {
		return nil, err
	}

	return
}

func (u *User) GenResetPassToken() (token string, err error) {

	o := orm.NewOrm()
	qs := o.QueryTable(u)

	for {
		token = common.RandStr(50, "alphanum")
		u.ResetPasswordToken = token

		if exist := qs.Filter("reset_password_token", token).Exist(); !exist {
			break
		}
	}

	t := time.Now()
	u.ResetPasswordSentAt = t
	_ , err = qs.Filter("id", u.Id).Update(orm.Params{"reset_password_token": token, "reset_password_sent_at": t,})

	return
}

func (u *User) ClearResetPassToken() error {

	o := orm.NewOrm()
	qs := o.QueryTable(u)

	t := time.Now()
	u.ResetPasswordSentAt = t
	_ , err := qs.Filter("id", u.Id).Update(orm.Params{"reset_password_token": "", "reset_password_sent_at": t,})

	return err
}

func (u *User) UpdatePassword(password string) (err error) {

	o := orm.NewOrm()

	if len(password) >= 6 {
		u.EncryptedPassword = common.Pwdhash(password)
	}

	o.Update(u, "encrypted_password")

	return
}

func UserGetByEmail(email string) (user *User, err error) {

	o := orm.NewOrm()
	user = new(User)
	user.Email = email
	if err = o.Read(user, "Email"); err != nil {
		return
	}

	if err = user.GetHistory(); err != nil {
		return
	}

	err = user.LoadRelated()

	return
}

func UserGetByResetPassToken(token string) (user *User, err error) {

	o := orm.NewOrm()

	if utf8.RuneCountInString(token) > 255 {
		return
	}

	user = new(User)
	user.ResetPasswordToken = token
	if err = o.QueryTable(user).Filter("reset_password_token", token).One(user); err != nil {
		return
	}

	t := time.Now()
	sub := t.Sub(user.ResetPasswordSentAt.Add(time.Hour * 24)).String()
	if !strings.Contains(sub, "-") {
		err = errors.New("max 24 hour")
	}

	return
}