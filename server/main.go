package main

import (
	"PIM_Server/config"
	routers "PIM_Server/routes"
	"PIM_Server/routes/auth"
	"PIM_Server/routes/chat"
	"PIM_Server/service"
	"flag"
	"log"
)

func main() {
	confFile := flag.String("f", "../etc/conf.yaml", "配置文件路径")
	flag.Parse()

	config.Init(*confFile)
	log.SetFlags(0)

	log.Printf("Start im_server, listen: %s\n", config.AppConfig().ServerInfo.Listen)
	//log.Fatal(http.ListenAndServe(config.AppConfig().ServerInfo.Listen, nil))

	service.StartClientMgr()

	// 加载多个APP的路由配置。有新增路由在此处注册
	routers.Register(service.Websocket)
	// auth相关通知
	routers.Register(auth.Routers, chat.Routers)

	// 初始化路由
	r := routers.Init()

	if err := r.Run(config.AppConfig().ServerInfo.Listen); err != nil {
		log.Fatalf("startup service failed, err:%v\n\n", err)
	}
}
