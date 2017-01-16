package models

type UserMeta struct {
	Id		int64		`orm:"pk;auto" json:"id"`
	User		*User		`orm:"rel(fk)" json:"-"`
	Key		string		`orm:"size(255)" valid:"MaxSize(255)" json:"key"`
	Value		string		`orm:"size(255)" valid:"MaxSize(255)" json:"value"`
}
