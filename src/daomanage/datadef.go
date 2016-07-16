package daomanage

import (
	"gopkg.in/mgo.v2/bson"
)

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
}

type User struct {
	UID          int64   `bson:"uid" json:"uid"`
	Username     string  `bson:"username" json:"username"`
	Password     string  `bson:"password" json:"password"`
	Nickname     string  `bson:"nickname" json:"nickname"`
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
	Index  string `bson:"_index" json:"_index"`
	PID    int64  `bson:"pid" json:"pid"`
	Status string `bson:"status" json:"status"`
	Code   string `bson:"code" json:"code"`
	Memory string `bson:"memory" json:"memory"`
	Time   string `bson:"time" json:"time"`
	Lang   string `bson:"lang" json:"lang"`
}
