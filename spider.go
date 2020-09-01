package nineone

import (
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"nineone/db"
	"reflect"
	"strconv"
	"time"
)

var listUrl = "https://api.zhaiclub.com/source/source_list"
var day = flag.Int("day", 5, "")
var fp = flag.Int("fp", 30, "")

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
	Image       string
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
	page := 0
	limit := 50
	errNo := 0
	//tu := timeUtil{time.Now()}
	//loc, _ := time.LoadLocation("Asia/Shanghai") //设置时区
	//zeroT := tu.ZeroTime().AddDate(0, 0, -*day)
	over := false
	for {


		res, err := s.fetchList(page, limit)
		if err != nil {
			if errNo >= 5 {
				log.Error("重试尝过3次以上", err)
			}
			log.Error(err)
			errNo += 1
			continue
		}
		fmt.Println(res.Data.List)

		errNo = 0
		page += 1

		data := []interface{}{}
		existNo := 0
		for _, v := range res.Data.List {
			//tt, _ := time.ParseInLocation("2006-01-02", v.UpTime, loc)
			//if tt.Before(zeroT) && !zeroT.Equal(tt) {
			if page > *fp {
				over = true
				break
			}
			exist := ! db.DB().Select("id").
				Where("vid = ?", v.Vid).First(&db.VideoList{}).RecordNotFound()
			if exist {
				existNo += 1
				continue
			}
			m := db.VideoList{}
			CopyStruct(&v, &m)
			data = append(data, m)
		}
		if over {
			log.Info("fetch over")
			return
		}

		log.Infof("存在 %d, 需要插入 %d", existNo, len(data))
		if len(data) > 0 {
			n, err := db.BatchInsert(data)
			log.Infof("插入了 %d, err: %v", n, err)
		}

		time.Sleep(time.Second * 5)
	}
}

type fetchDetailRes struct {
	Code int
	Msg string
	Data string
}
func (s *Spider) FetchDetail(uri string) (string, error) {
	resJson := fetchDetailRes{}

	md5Val := fmt.Sprintf("%x", md5.Sum([]byte(uri)))
	if n := db.Redis().Exists(md5Val).Val(); n > 0 {
		return db.Redis().Get(md5Val).String(), nil
	}

	uri = fmt.Sprintf("%s?key=%s&act=url", uri, "TJw92fDnLsChYvkX")
	resp, err := RequestBySocket5(uri)
	//req, _ := BuildGetReq(uri, "GET", item)
	//res, err := req.Do()
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(resp, &resJson); err != nil {
		return "", err
	}
	db.Redis().Set(md5Val, resJson.Data, time.Hour)
	return resJson.Data, nil
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
