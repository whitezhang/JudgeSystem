package auth

import (
	"code.google.com/p/gcfg"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	AuthCode_Unknown      = -1
	AuthCode_OK           = 0
	AUthCode_QpsOverload  = 1
	AuthCode_Unauthorized = 2
)

type QpsLimit struct {
	maxQps    int
	curQps    int
	timeStamp int64
}

type Manager struct {
	authInfo map[string]*QpsLimit
	mutex    *sync.Mutex
}

func NewManager(cfgFile string) (man *Manager, err error) {
	man = &Manager{}
	man.authInfo = make(map[string]*QpsLimit)
	man.mutex = new(sync.Mutex)
	err = man.initConf(cfgFile)
	return
}

func (man *Manager) initConf(cfgFile string) (err error) {
	var cfg authConfig
	err = gcfg.ReadFileInto(&cfg, cfgFile)
	if err != nil {
		return
	}
	log.Println("Auth key: ")
	for _, auth := range cfg.Auth {
		k := auth.IpAddr
		fmt.Println("\t", k)
		v := new(QpsLimit)
		v.maxQps = auth.MaxQps
		v.curQps = 0
		v.timeStamp = time.Now().Unix()
		man.authInfo[k] = v
	}
	return
}

func (man *Manager) DoAuth(ipAddr string) (authCode int, err error) {
	authCode = AuthCode_Unknown
	err = nil

	authKey := ipAddr
	man.mutex.Lock()
	defer man.mutex.Unlock()
	authValue, ok := man.authInfo[authKey]
	if ok {
		now := time.Now().Unix()
		if now != authValue.timeStamp {
			authValue.curQps = 0
			authValue.timeStamp = now
		}
		authValue.curQps += 1
		if authValue.maxQps < authValue.curQps {
			authCode = AUthCode_QpsOverload
		} else {
			authCode = AuthCode_OK
		}
	} else {
		authCode = AuthCode_Unauthorized
	}
	return
}
