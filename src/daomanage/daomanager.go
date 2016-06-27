package daomanage

import (
	"code.google.com/p/gcfg"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"
)

var (
	hostName string
	port     string
	dbName   string
)

type Manager struct {
	serviceMap map[string]string
	mutex      *sync.Mutex
}

func NewManager(cfgFile string) (man *Manager, err error) {
	man = &Manager{}
	man.serviceMap = make(map[string]string)
	man.mutex = new(sync.Mutex)
	err = man.initConf(cfgFile)
	return
}

func (man *Manager) initConf(cfgFile string) (err error) {
	var cfg daoConfig
	err = gcfg.ReadFileInto(&cfg, cfgFile)
	if err != nil {
		return
	}
	hostName = cfg.Dao.HostName
	port = cfg.Dao.Port
	dbName = cfg.Dao.DBName
	log.Printf("Dao conf: %s:%s:%s\n", hostName, port, dbName)
	return
}

func GetContest(cid string) (contestInfo Contest, err error) {
	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return Contest{}, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(dbName).C("contest")

	err = collection.Find(bson.M{"cid": cid}).One(&contestInfo)
	if err != nil {
		log.Printf("No Contest named: %s\n", cid)
		return Contest{}, err
	}
	return contestInfo, nil
}

func GetContestProblems(cid string) (contestInfo []ContestProblem, err error) {
	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return []ContestProblem{}, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(dbName).C("contestproblemlist")

	err = collection.Find(bson.M{"cid": cid}).All(&contestInfo)
	if err != nil {
		log.Printf("No Contest named: %s\n", cid)
		return []ContestProblem{}, err
	}
	return contestInfo, nil
}

func GetProblemInfo(pid string) (problemInfo Problem, err error) {
	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return Problem{}, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(dbName).C("problem")

	err = collection.Find(bson.M{"pid": pid}).One(&problemInfo)
	if err != nil {
		log.Printf("No Problem named: %s\n", pid)
		return Problem{}, err
	}
	return
}

func GetUserInfo(key, value string) (userInfo User, err error) {
	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return User{}, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(dbName).C("user")

	err = collection.Find(bson.M{key: value}).One(&userInfo)
	// log.Println(userInfo)
	if err != nil {
		log.Printf("No User %s: %s\n", key, value)
		return User{}, err
	}
	return
}
