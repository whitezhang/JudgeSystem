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

type ProblemInfo struct {
	PID         int64  `bson:"pid" json:"pid"`
	ProblemName string `bson:"title" json:"title"`
	Solved      int64  `bson:"solved" json:"solved"`
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

func InsertSubmitQueue(pid int64, code string, lang string) (err error) {
	var exSbtQue []ExSubmitQueue
	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	sbtCollection := session.DB(dbName).C("submitqueue")
	rsCollection := session.DB(dbName).C("runtimestatus")

	err = sbtCollection.Insert(&SubmitQueue{pid, code, lang})
	if err != nil {
		log.Println("Error: Failed in submition")
		return err
	}
	err = sbtCollection.Find(bson.M{"pid": pid, "code": code, "lang": lang}).All(&exSbtQue)
	if err != nil {
		log.Println("Error: Failed in submition")
		return err
	}
	for _, submition := range exSbtQue {
		err = rsCollection.Insert(&RuntimeStatus{submition.ID.Hex(), pid, "Pending", code, "Pending", "Pending", lang})
		if err != nil {
			log.Println("Error: Failed in submition")
			continue
		}
	}
	return nil
}

func GetContest(cid int64) (contestInfo Contest, err error) {
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

func GetContestProblems(cid int64) (contestInfo []ContestProblem, err error) {
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

func GetProblemInRange(startIndex, endIndex int64) (problemInfo []ProblemInfo, err error) {
	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return []ProblemInfo{}, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(dbName).C("problem")

	err = collection.Find(bson.M{"pid": bson.M{"$gte": startIndex, "$lte": endIndex}}).Select(bson.M{"pid": 1, "title": 1, "solved": 1}).All(&problemInfo)
	if err != nil {
		log.Printf("No Problem Indexing: from %d to %d\n", startIndex, endIndex)
		return []ProblemInfo{}, err
	}
	return
}

func GetProblemInfo(pid int64) (problemInfo Problem, err error) {
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
	if err != nil {
		log.Printf("No User %s: %s\n", key, value)
		return User{}, err
	}
	return
}
