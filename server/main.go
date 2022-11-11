package main

import (
	"PIM_Server/config"
	"PIM_Server/service"
	"flag"
	"log"
	"net/http"
)

func main() {
	confFile := flag.String("f", "../etc/conf.yaml", "配置文件路径")
	flag.Parse()

	config.Init(*confFile)
	log.SetFlags(0)
	http.HandleFunc("/ws", service.WS)

	log.Printf("Start im_server success, listen: %s\n", config.AppConfig().ServerInfo.Listen)
	log.Fatal(http.ListenAndServe(config.AppConfig().ServerInfo.Listen, nil))
}
