package main

import (
	"fmt"
	"time"
)

func main() {
	loc, _ := time.LoadLocation("Asia/Shanghai")                        //设置时区
	now := time.Now()
	d := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	d = d.AddDate(0,0,-2)
	s := "2020-08-28"
	tt, _ := time.ParseInLocation("2006-01-02", s, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	fmt.Println(tt, d)
	fmt.Println(tt.Before(d))
}

