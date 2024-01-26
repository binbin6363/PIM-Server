package main

import (
	"PIM_Server/config"
	"PIM_Server/log"
	"PIM_Server/plugins"
	"PIM_Server/service"
	"flag"
	"os"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func main() {
	confFile := flag.String("f", "../etc/conf.yaml", "配置文件路径")
	flag.Parse()
	log.InitLogger(config.AppConfig().LogInfo.Path,
		config.AppConfig().LogInfo.MaxSize,
		config.AppConfig().LogInfo.MaxBackUps,
		config.AppConfig().LogInfo.MaxAge,
		config.AppConfig().LogInfo.Level,
		config.AppConfig().LogInfo.CallerSkip)

	config.Init(*confFile)

	log.Infof("Start im_server, listen: %s", config.AppConfig().ServerInfo.Listen)

	service.Init()
	service.StartClientMgr()

	r := plugins.Init(serviceName)
	if err := r.Run(config.AppConfig().ServerInfo.Listen); err != nil {
		log.Fatalf("startup service failed, err:%v", err)
	}
}
