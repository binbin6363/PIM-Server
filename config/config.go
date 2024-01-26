package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ServerInfo struct {
	Name         string `yaml:"name"`
	Listen       string `yaml:"listen"`
	Timeout      int    `yaml:"timeout"`
	Secret       string `yaml:"secret"`
	TokenExpire  int    `yaml:"token_expire"`
	DataCenterId int64  `yaml:"data_center_id"`
	WorkerId     int64  `yaml:"worker_id"`
	DebugReqRsp  bool   `yaml:"debug_req_rsp"`
	Mode         string `yaml:"mode"`
}

type LogInfo struct {
	Path       string `yaml:"path"`
	Level      int    `yaml:"level"`
	MaxSize    int    `yaml:"max_size"` // mb
	MaxAge     int    `yaml:"max_age"`  // day
	MaxBackUps int    `yaml:"max_backups"`
	CallerSkip int    `yaml:"caller_skip"`
}

type DBInfo struct {
	// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	Dsn          string `yaml:"dsn"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxLifeTime  int    `yaml:"max_life_time"` // 单位秒
}

type ServerCfg struct {
	ServerInfo *ServerInfo `yaml:"server"`
	DBInfo     *DBInfo     `yaml:"db"`
	LogInfo    *LogInfo    `yaml:"log"`
}

// 配置实例
var cfg = &ServerCfg{}

// AppConfig 获取配置单例
func AppConfig() *ServerCfg {
	return cfg
}

// Init 初始化配置
func Init(file string) {
	configFile, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("load conf fail, path:%s, err:%v", file, err)
	}

	if err = yaml.Unmarshal(configFile, cfg); err != nil {
		log.Fatalf("Unmarshal conf fail, err:%v", err)
	}

	log.Printf("load conf ok, path:%s, conf:%v", file, string(configFile))
}
