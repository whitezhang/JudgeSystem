package context

import (
	"config"
	"daomanage"
	"errors"
	"flag"
	"fmt"
	// "log"
	//"service/auth"
)

var (
	SvrCtx ServerContext
)

type ServerContext struct {
	SvrCfg config.ServerConfig
	//AuthMan *auth.Manager
	DaoMan *daomanage.Manager
}

func InitServerContext() error {
	SvrCtx = ServerContext{}

	//var authCfgFile = flag.String("authconf", "./conf/auth.conf", "the config file for auth module.")
	var serverCfgFile = flag.String("serverconf", "./conf/unserver.conf", "the config file for server module.")
	var daoCfgFile = flag.String("daoconf", "./conf/dao.conf", "the config file for dao module.")
	flag.Parse()

	if *serverCfgFile == "" {
		return errors.New("Empty Auth Config")
	}
	if err := config.InitConfig(&SvrCtx.SvrCfg, *serverCfgFile); err != nil {
		return err
	}
	/*
		if *authCfgFile == "" {
			return errors.New("Empty Auth Config")
		}
		if authMan, err := auth.NewManager(*authCfgFile); err == nil {
			SvrCtx.AuthMan = authMan
		} else {
			return fmt.Errorf("Failed to init auth manager, error: %s", err)
		}
	*/

	if *daoCfgFile == "" {
		return errors.New("Empty Dao Config")
	}
	if daoMan, err := daomanage.NewManager(*daoCfgFile); err == nil {
		SvrCtx.DaoMan = daoMan
	} else {
		return fmt.Errorf("Failed to init dao manager, error: %s", err)
	}

	return nil
}
