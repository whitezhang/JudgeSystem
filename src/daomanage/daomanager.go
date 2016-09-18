package daomanage

import (
	"code.google.com/p/gcfg"
	ctx "context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math"
	"sync"
)

var (
	PlbPerPage  = ctx.SvrCtx.SvrCfg.WebInfo.ProblemPerPage
	StatPerPage = ctx.SvrCtx.SvrCfg.WebInfo.StatusPerPage
	CstPerPage  = ctx.SvrCtx.SvrCfg.WebInfo.ContestPerPage
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

/*
type ContestInfoSet struct {
	PageCount   int       `bson:"pagecount" json:"pagecount"`
	CstInfoList []Contest `bson:"cstinfolist" json:"cstinfolist"`
}

type StatInfo struct {
}

type StatusInfoSet struct {
	PageCount    int             `bson:"pagecount" json:"pagecount"`
	StatInfoList []RuntimeStatus `bson:"statinfolist" json:"statinfolist"`
}

type ProblemInfoSet struct {
	PageCount   int           `bson:"pagecount" json:"pagecount"`
	ProInfoList []ProblemInfo `bson:"proinfolist" json:"proinfolist"`
}

type ProblemInfo struct {
	PID         int64  `bson:"pid" json:"pid"`
	ProblemName string `bson:"title" json:"title"`
	Solved      int64  `bson:"solved" json:"solved"`
}
*/

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

func InsertRegister(email, username, password, challenger string) (err error) {
	var userinfo User
	var ischallenger bool
	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(dbName).C("user")
	err = collection.Find(bson.M{"username": username}).One(&userinfo)
	if err == nil {
		log.Printf("The username %s has been registered", username)
		return err
	}
	if challenger == "1" {
		ischallenger = true
	} else {
		ischallenger = false
	}
	err = collection.Insert(&User{email, username, password, ischallenger, 0.0, "traveller"})
	if err != nil {
		log.Printf("Error: Failed in register")
		return err
	}
	return nil
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

func GetStatusInRange(startIndex, endIndex int64) (statusInfoSet StatusInfoSet, err error) {
	var statInfo []RuntimeStatus
	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return StatusInfoSet{}, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(dbName).C("runtimestatus")

	cnt, err := collection.Count()
	cnt = int(math.Ceil(float64(cnt) / StatPerPage))
	if err != nil {
		log.Printf("Get runtimestatus count error\n")
		return StatusInfoSet{}, err
	}

	err = collection.Find(bson.M{"pid": bson.M{"$gte": startIndex, "$lte": endIndex}}).Select(bson.M{"pid": 1, "status": 1, "memory": 1, "time": 1, "lang": 1}).All(&statInfo)
	if err == nil {
		log.Printf("No Runtimestatus indexing: from: %d to %d\n", startIndex, endIndex)
		statusInfoSet.PageCount = cnt
		statusInfoSet.StatInfoList = statInfo
		return statusInfoSet, err
	}

	return
}

func GetContestInRange(startIndex, endIndex int64) (contestInfoSet ContestInfoSet, err error) {
	var contestInfo []Contest

	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return ContestInfoSet{}, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(dbName).C("contest")

	cnt, err := collection.Count()
	cnt = int(math.Ceil(float64(cnt) / CstPerPage))
	if err != nil {
		log.Printf("Get contest count error\n")
		return ContestInfoSet{}, err
	}

	err = collection.Find(bson.M{"cid": bson.M{"$gte": startIndex, "$lte": endIndex}}).All(&contestInfo)
	if err == nil {
		log.Printf("No Contest indexing: from: %d to %d\n", startIndex, endIndex)
		contestInfoSet.PageCount = cnt
		contestInfoSet.CstInfoList = contestInfo
		return contestInfoSet, err
	}

	return
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

//func GetProblemInRange(startIndex, endIndex int64) (problemInfo []ProblemInfo, err error) {
func GetProblemInRange(startIndex, endIndex int64) (problemInfoSet ProblemInfoSet, err error) {
	var problemInfo []ProblemInfo

	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return ProblemInfoSet{}, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(dbName).C("problem")

	cnt, err := collection.Count()
	cnt = int(math.Ceil(float64(cnt) / PlbPerPage))
	if err != nil {
		log.Printf("Get problem count error\n")
		return ProblemInfoSet{}, err
	}

	err = collection.Find(bson.M{"pid": bson.M{"$gte": startIndex, "$lte": endIndex}}).Select(bson.M{"pid": 1, "title": 1, "solved": 1}).All(&problemInfo)
	if err == nil {
		log.Printf("No Problem Indexing: from %d to %d\n", startIndex, endIndex)
		problemInfoSet.PageCount = cnt
		problemInfoSet.ProInfoList = problemInfo
		return problemInfoSet, err
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
