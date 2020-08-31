package nineone

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	TimeLayout = "2006-01-02 15:04:05"
)

func GetBeforeDaysTime(d time.Time, day int) time.Time {
	return d.AddDate(0, 0, -day)
}

func GetFirstDateOfYear(d time.Time) time.Time {
	return time.Date(d.Year(), 1, 1, 0, 0, 0, 0, d.Location())
}

type timeUtil struct {
	time time.Time
}

func (t *timeUtil) FromHourString(date string) (time.Time, bool) {
	tmp := strings.Split(date, ":")
	if len(tmp) != 2 {
		return t.time, false
	}
	hour, err1 := strconv.Atoi(tmp[0])
	minute, err2 := strconv.Atoi(tmp[1])
	if err1 != nil || err2 != nil {
		return t.time, false
	}

	return time.Date(t.time.Year(), t.time.Month(), t.time.Day(),
		hour, minute, 0, 0, t.time.Location()), true
}

func (t *timeUtil) Time() time.Time {
	return t.time
}

//获取某一天的0点时间
func (t *timeUtil) ZeroTime() time.Time {
	return time.Date(t.time.Year(), t.time.Month(), t.time.Day(), 0, 0, 0, 0, t.time.Location())
}
func (t *timeUtil) EndTime() time.Time {
	return t.ZeroTime().AddDate(0,0,1)
}

func (t *timeUtil) BeforeDaysTime(day int) time.Time {
	return t.time.AddDate(0, 0, -day)
}

func (t *timeUtil) BeforeDaysTimeAndZero(day int) time.Time {
	return t.ZeroTime().AddDate(0, 0, day)
}

func (t *timeUtil) FirstDateOfYear() time.Time {
	return time.Date(t.time.Year(), 1, 1, 0, 0, 0, 0, t.time.Location())
}

func NewTimeUtil(time time.Time) *timeUtil {
	return &timeUtil{time: time}
}

func ParseTime(date string) (time.Time, error) {
	return time.ParseInLocation(TimeLayout, date, Location)
}

// 获取毫秒数
func GetTimestampByMs() int64{
	return time.Now().UnixNano() / 1e6
}

var Location *time.Location

func SetLocation() {
	var err error
	Location, err = time.LoadLocation("Local")
	if err != nil {
		logrus.Fatal(err)
	}
}

func FilterTime(date string) time.Time {
	var startTime time.Time
	timeUtil := NewTimeUtil(time.Now())
	switch date {
	case "7days":
		startTime = timeUtil.BeforeDaysTime(7)
	case "30days":
		startTime = timeUtil.BeforeDaysTime(30)
	case "all":
		startTime = timeUtil.FirstDateOfYear()
	default:
		// 默认今天
		startTime = timeUtil.ZeroTime()
	}
	return startTime
}


