package main

import (
	"encoding/json"
	"flag"
	"github.com/sirupsen/logrus"
	"net/http"
	"nineone"
	"nineone/config"
	"nineone/db"
	"strconv"
)

func toJson(data interface{}) []byte{
	b, _ := json.Marshal(data)
	return b
}

func videoList(w http.ResponseWriter, req *http.Request) {
	page := req.URL.Query().Get("page")
	iPage, _ := strconv.Atoi(page)
	limit := req.URL.Query().Get("per_page")
	keyword := req.URL.Query().Get("keyword")
	author := req.URL.Query().Get("author")
	iLimit, _ := strconv.Atoi(limit)
	if iLimit == 0 {
		iLimit = 10
	}

	iPage = (iPage - 1) * iLimit
	var data []db.VideoList
	m := db.DB().Offset(iPage)
	if keyword != "" {
		m = m.Where("title like ?", "%"+keyword+"%")
	}
	if author != "" {
		m = m.Where("author = ?", author)
	}
	m.Limit(iLimit).Order("up_time desc, id desc").Find(&data)
	w.Write(toJson(data))
	return
}

func videoUrl(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	s := nineone.NewSpider()
	ret, err := s.FetchDetail(id)
	if err == nil {
		b, _ := json.Marshal(ret)
		w.Write(b)
	} else {
		w.Write([]byte("error"))
		logrus.Error(err)
	}
	return
}

var configFile = flag.String("config", "/Users/zxy/project/go_projects/nineone/.env.yaml", "配置文件路径")

func main() {
	flag.Parse()
	config.InitConfig(*configFile)
	db.InitDB()
	http.HandleFunc("/video-list", videoList)
	http.HandleFunc("/video-url", videoUrl)
	http.ListenAndServe(":8090", nil)
}

