package null

import (
	"database/sql/driver"
	"strconv"
)

type Int64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

func NewInt64(value interface{}) (i Int64) {
	i = Int64{}
	i.Scan(value)
	return
}

func (n *Int64) Scan(value interface{}) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Int64, value)
}

func (n Int64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}


func (n Int64) String() (value string) {
	if !n.Valid {
		value = "null"
		return
	}
	value = strconv.FormatInt(n.Int64, 10)
	return
}

func (n Int64) MarshalJSON() ([]byte, error) {
	return []byte(n.String()), nil
}

func (n *Int64) UnmarshalJSON(data []byte) error {
	i64, err := strconv.ParseInt(string(data), 10, 0)
	if err != nil {
		return nil
	}
	n.Valid = true
	return n.Scan(i64)
}