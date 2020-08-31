package main

import (
	"flag"
	"nineone"
	"nineone/config"
	"nineone/db"
)

var configFile = flag.String("config", "/Users/zxy/project/go_projects/nineone/.env.yaml", "配置文件路径")

func main() {
	flag.Parse()

	config.InitConfig(*configFile)
	db.InitDB()

	nineone.SetLocation()
	spider := nineone.NewSpider()
	spider.Run()
}
