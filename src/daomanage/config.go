package daomanage

import ()

type daoCfg struct {
	HostName string
	Port     string
	DBName   string
}

type daoConfig struct {
	Dao daoCfg
}
