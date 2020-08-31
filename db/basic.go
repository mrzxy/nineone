package db

import (
	"database/sql/driver"
	"fmt"
	"math"
	"time"
)

type BasicModel struct {
	ID int `gorm:"primary_key" json:"id"`
}

type Timestamp struct {
	CreatedAt JSONTime  `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime" json:"-"`
}

type JSONTime struct {
	time.Time
}

// Value insert timestamp into mysql need this function.
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *JSONTime) DiffFormat() string {
	ct := time.Now().Unix() - t.Time.Unix()
	if ct < 0 {
		return "刚刚"
	} else if ct < 60 {
		return fmt.Sprintf("%d秒前", ct)
	} else if ct < 3600 {
		return fmt.Sprintf("%d分钟前", int64(math.Floor(float64(ct/60))))
	} else if ct < 3600*12 {
		return fmt.Sprintf("%d小时前", int64(math.Floor(float64(ct/3600))))
	} else if ct < 365*86400 {
		return t.Time.Format("01-02 15:04")
	} else {
		return t.Time.Format("2006-01-02 15:04")
	}
}
