package daomanage

import (
	"code.google.com/p/gcfg"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

var (
	plbPerPage  int
	statPerPage int
	cstPerPage  int
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

func NewManager(cfgFile string, plbNum, statNum, cstNum int) (man *Manager, err error) {
	man = &Manager{}
	man.serviceMap = make(map[string]string)
	man.mutex = new(sync.Mutex)
	err = man.initConf(cfgFile, plbNum, statNum, cstNum)
	return
}

func (man *Manager) initConf(cfgFile string, plbNum, statNum, cstNum int) (err error) {
	var cfg daoConfig
	err = gcfg.ReadFileInto(&cfg, cfgFile)
	if err != nil {
		return
	}
	hostName = cfg.Dao.HostName
	port = cfg.Dao.Port
	dbName = cfg.Dao.DBName
	plbPerPage = plbNum
	statPerPage = statNum
	cstPerPage = cstNum
	log.Printf("Dao conf: %s:%s:%s\n", hostName, port, dbName)
	return
}

func GetNextID(idtype string) (index int, err error) {
	type Index struct {
		index int
	}
	var idx Index

	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return -1, err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(dbName).C("idmanager")

	err = collection.Update(bson.M{"name": idtype}, bson.M{"$inc": bson.M{"index": 1}})
	//err = collection.Update(bson.M{"name": idtype}, bson.M{"$inc": bson.M{"index": 1}})
	if err != nil {
		log.Println("Failed in IDMan:update")
		return -1, err
	}

	err = collection.Find(bson.M{"name": idtype}).One(&idx)
	if err != nil {
		log.Println("Failed in IDMan:get")
		return -1, err
	}
	log.Println(idx.index)
	return idx.index, nil
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
		log.Printf("The username %s has been registed", username)
		return errors.New("The username has been registed")
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

func InsertSubmitQueue(pid int, code string, lang string, author string) (err error) {
	var exSbtQue []ExSubmitQueue
	sbmtime := time.Now().Format("2006-01-02 15:04:05")
	codelen := strconv.Itoa(len(code))

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
		rid, err := GetNextID("rid")
		if err != nil {
			log.Println("Error: Failed in submition:getnextRid")
			return err
		}
		err = rsCollection.Insert(&RuntimeStatus{submition.ID.Hex(), rid, pid, sbmtime, "Pending", code, "Pending", "Pending", codelen, lang, author})
		if err != nil {
			log.Println("Error: Failed in submition")
			return err
		}
	}
	return nil
}

func GetStatusInRange(startIndex, endIndex int) (statusInfoSet StatusInfoSet, err error) {
	var statInfo []RuntimeStatusInfo
	//var statInfoList []RuntimeStatusInfo

	//_startIndex := int(startIndex)
	//_endIndex := int(endIndex)

	session, err := mgo.Dial(hostName)
	if err != nil {
		log.Println("Connect MongoDB failed")
		return StatusInfoSet{}, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	collection := session.DB(dbName).C("runtimestatus")

	cnt, err := collection.Count()
	cnt = int(math.Ceil(float64(cnt) / float64(statPerPage)))
	if err != nil {
		log.Printf("Get runtimestatus count error\n")
		return StatusInfoSet{}, err
	}

	/*
		iter := collection.Find(nil).Iter()
		index := 0
		for iter.Next(&statInfo) {
			if index >= _startIndex && index < _endIndex {
				statInfoList = append(statInfoList, statInfo)
			}
			index += 1
		}
	*/

	err = collection.Find(bson.M{"rid": bson.M{"$gte": startIndex, "$lte": endIndex}}).Select(bson.M{"rid": 1, "pid": 1, "sbmtime": 1, "status": 1, "memory": 1, "time": 1, "lang": 1, "codelen": 1, "author": 1}).All(&statInfo)
	if err != nil {
		log.Printf("No Runtimestatus indexing: from: %d to %d\n", startIndex, endIndex)
		statusInfoSet.PageCount = cnt
		statusInfoSet.StatInfoList = nil
		return
	}
	//err = collection.Find(bson.M{}).Select(bson.M{"pid": 1, "sbmtime": 1, "status": 1, "memory": 1, "time": 1, "lang": 1, "codelen": 1, "author": 1}).All(&statInfo)
	statusInfoSet.PageCount = cnt
	statusInfoSet.StatInfoList = statInfo
	return
}

func GetContestInRange(startIndex, endIndex int) (contestInfoSet ContestInfoSet, err error) {
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
	cnt = int(math.Ceil(float64(cnt) / float64(cstPerPage)))
	if err != nil {
		log.Printf("Get contest count error\n")
		return ContestInfoSet{}, err
	}

	err = collection.Find(bson.M{"cid": bson.M{"$gte": startIndex, "$lt": endIndex}}).All(&contestInfo)
	if err == nil {
		log.Printf("No Contest indexing: from: %d to %d\n", startIndex, endIndex)
		contestInfoSet.PageCount = cnt
		contestInfoSet.CstInfoList = contestInfo
		return contestInfoSet, err
	}

	return
}

func GetContestProblems(cid int) (contestInfo []ContestProblem, err error) {
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

//func GetProblemInRange(startIndex, endIndex int) (problemInfo []ProblemInfo, err error) {
func GetProblemInRange(startIndex, endIndex int) (problemInfoSet ProblemInfoSet, err error) {
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
	cnt = int(math.Ceil(float64(cnt) / float64(plbPerPage)))
	if err != nil {
		log.Printf("Get problem count error\n")
		return ProblemInfoSet{}, err
	}

	err = collection.Find(bson.M{"pid": bson.M{"$gte": startIndex, "$lt": endIndex}}).Select(bson.M{"pid": 1, "title": 1, "solved": 1, "author": 1}).All(&problemInfo)
	if err == nil {
		log.Printf("No Problem Indexing: from %d to %d\n", startIndex, endIndex)
		problemInfoSet.PageCount = cnt
		problemInfoSet.ProInfoList = problemInfo
		return problemInfoSet, err
	}
	return
}

func GetProblemInfo(pid int) (problemInfo Problem, err error) {
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
