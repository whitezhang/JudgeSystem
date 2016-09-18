package config

import (
	"code.google.com/p/gcfg"
)

type Log4goCfg struct {
	ProgName    string
	Dir         string
	Level       string
	BackupCount int
	When        string
	HasStdout   bool
}

type WebInfoCfg struct {
	ProblemPerPage int64
	StatusPerPage  int64
	ContestPerPage int64
}

type ServerCfg struct {
	Port   int
	NumCPU int
}

type ServerConfig struct {
	Server  ServerCfg
	Log4go  Log4goCfg
	WebInfo WebInfoCfg
}

func InitConfig(cfg *ServerConfig, cfgFile string) error {
	return gcfg.ReadFileInto(cfg, cfgFile)
}
