package nineone

import (
	"flag"
	"net/url"
	"nineone/db"
	"reflect"
	"strconv"
	"time"
)

var listUrl = "https://api.zhaiclub.com/source/source_list"
var day = flag.Int("day", 5, "")

type listResult struct {
	Status int
	Msg    string
	Count  string
	Data   ListResultData
}

type ListResultData struct {
	List []listResultItem
}

type listResultItem struct {
	Viewkey   string
	Img       string
	Author    string
	UpTime    string `json:"up_time"`
	Title     string
	Vid       string
	Duration  string
	View      int
	Favorites int
	Comment   int
	Integral  int
	VideoUrl  string `json:"video_url"`
}

type Spider struct {
}

func (s *Spider) fetchList(page int, limit int) (listResult, error) {
	resJson := listResult{}
	item := url.Values{}
	item.Set("page", strconv.Itoa(page))
	item.Set("limit", strconv.Itoa(limit))

	req, _ := BuildGetReq(listUrl, "GET", item)
	res, err := req.Do()
	if err != nil {
		return resJson, err
	}
	if err := res.Body.FromJsonTo(&resJson); err != nil {
		return resJson, err
	}
	return resJson, nil
}
func (s *Spider) Run() {
	log := NewLogrus("fetchList")
	page := 1
	limit := 50
	errNo := 0
	tu := timeUtil{time.Now()}
	loc, _ := time.LoadLocation("Asia/Shanghai") //设置时区
	zeroT := tu.ZeroTime().AddDate(0, 0, -*day)
	over := false
	for {
		if over {
			log.Info("fetch over")
			return
		}

		res, err := s.fetchList(page, limit)
		if err != nil {
			if errNo >= 3 {
				log.Error("重试尝过3次以上", err)
			}
			log.Error(err)
			errNo += 1
			continue
		}

		errNo = 0
		page += 1

		data := []interface{}{}
		existNo := 0
		for _, v := range res.Data.List {
			tt, _ := time.ParseInLocation("2006-01-02", v.UpTime, loc)
			if tt.Before(zeroT) && !zeroT.Equal(tt) {
				over = true
				break
			}
			exist := ! db.DB().Select("id").
				Where("vid = ?", v.Vid).First(db.VideoList{}).RecordNotFound()
			if exist {
				existNo += 1
			}
			m := db.VideoList{}
			CopyStruct(&v, &m)
			data = append(data, m)
		}
		log.Infof("存在 %d, 需要插入 %d", existNo, len(data))
		if len(data) > 0 {
			n, err := db.BatchInsert(data)
			log.Infof("插入了 %d, err: %v", n, err)
		}
		time.Sleep(time.Second * 5)
	}
}

func NewSpider() *Spider {
	return &Spider{}
}

func CopyStruct(src, dst interface{}) {
	sval := reflect.ValueOf(src).Elem()
	dval := reflect.ValueOf(dst).Elem()

	for i := 0; i < sval.NumField(); i++ {
		value := sval.Field(i)
		name := sval.Type().Field(i).Name

		dvalue := dval.FieldByName(name)
		if dvalue.IsValid() == false || dvalue.Type() != value.Type() {
			continue
		}
		dvalue.Set(value)
	}
}
