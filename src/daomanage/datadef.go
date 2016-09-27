package daomanage

import (
	"gopkg.in/mgo.v2/bson"
)

type IDManager struct {
	PID int64 `bson:"pid" json:"pid"`
	CID int64 `bson:"cid" json:"cid"`
	RID int64 `bson:"rid" json:"rid"`
}

type Problem struct {
	PID          int64  `bson:"pid" json:"pid"`
	Title        string `bson:"title" json:"title"`
	Description  string `bson:"description" json:"description"`
	Time         string `bson:"time" json:"time"`
	Memory       string `bson:"memory" json:"memory"`
	Input        string `bson:"input" json:"input"`
	Output       string `bson:"output" json:"output"`
	SimpleInput  string `bson:"simpleinput" json:"simpleinput"`
	SimpleOutput string `bson:"simpleoutput" json:"simpleoutput"`
	Solved       int64  `bson:"solved" json:"solved"`
	Display      bool   `bson:"display" json:"display"`
	Author       string `bson:"author" json:"author"`
}

type User struct {
	Email        string  `bson:"email" json:"email"`
	Username     string  `bson:"username" json:"username"`
	Password     string  `bson:"password" json:"password"`
	IsChallenger bool    `bson:"ischallenger" json:"ischallenger"`
	Score        float64 `bson:"score" json:"score"`
	Privilege    string  `bson:"privilege" json:"privilege"`
}

type Contest struct {
	CID int64 `bson:"cid" json:"cid"`
	// pid1;pid2;pid3
	ContestName string `bson:"contestname" json:"contestname"`
	ProblemList string `bson:"problemlist" json:"problemlist"`
	StartTime   int64  `bson:"starttime" json:"starttime"`
	EndTime     int64  `bson:"endtime" json:"endtime"`
}

type ContestProblem struct {
	CID    int64 `bson:"cid" json:"cid"`
	PID    int64 `bson:"pid" json:"pid"`
	Solved int64 `bson:"solved" json:"solved"`
	Score  int64 `bson:"score" json:"score"`
}

type SubmitQueue struct {
	PID  int64  `bson:"pid" json:"pid"`
	Code string `bson:"code" json:"code"`
	Lang string `bson:"lang" json:"lang"`
}

type ExSubmitQueue struct {
	ID   bson.ObjectId `bson:"_id" json:"_id"`
	PID  int64         `bson:"pid" json:"pid"`
	Code string        `bson:"code" json:"code"`
	Lang string        `bson:"lang" json:"lang"`
}

type RuntimeStatus struct {
	Index   string `bson:"_index" json:"_index"`
	RID     int64  `bson:"rid" json:"rid"`
	PID     int64  `bson:"pid" json:"pid"`
	SbmTime string `bson:"sbmtime" json:"sbmtime"`
	Status  string `bson:"status" json:"status"`
	Code    string `bson:"code" json:"code"`
	Memory  string `bson:"memory" json:"memory"`
	Time    string `bson:"time" json:"time"`
	CodeLen string `bson:"codelen" json:"codelen"`
	Lang    string `bson:"lang" json:"lang"`
	Author  string `bson:"author" json:"author"`
}

// Only for golang/javascript communication
type ContestInfoSet struct {
	PageCount   int       `bson:"pagecount" json:"pagecount"`
	CstInfoList []Contest `bson:"cstinfolist" json:"cstinfolist"`
}

type StatusInfoSet struct {
	PageCount    int                 `bson:"pagecount" json:"pagecount"`
	StatInfoList []RuntimeStatusInfo `bson:"statinfolist" json:"statinfolist"`
}

type RuntimeStatusInfo struct {
	RID     int64  `bson:"rid" json:"rid"`
	PID     int64  `bson:"pid" json:"pid"`
	SbmTime string `bson:"sbmtime" json:"sbmtime"`
	Status  string `bson:"status" json:"status"`
	Memory  string `bson:"memory" json:"memory"`
	Time    string `bson:"time" json:"time"`
	CodeLen string `bson:"codelen" json:"codelen"`
	Lang    string `bson:"lang" json:"lang"`
	Author  string `bson:"author" json:"author"`
}

type ProblemInfoSet struct {
	PageCount   int           `bson:"pagecount" json:"pagecount"`
	ProInfoList []ProblemInfo `bson:"proinfolist" json:"proinfolist"`
}

type ProblemInfo struct {
	PID         int64  `bson:"pid" json:"pid"`
	ProblemName string `bson:"title" json:"title"`
	Solved      int64  `bson:"solved" json:"solved"`
	Author      string `bson:"author" json:"author"`
}
