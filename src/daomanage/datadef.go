package daomanage

import ()

type Problem struct {
	PID          string
	Title        string `bson:"title" json:"title"`
	Description  string `bson:"description" json:"description"`
	Time         string `bson:"time" json:"time"`
	Memory       string `bson:"memory" json:"memory"`
	SimpleInput  string `bson:"simpleinput" json:"simpleinput"`
	SimpleOutput string `bson:"simpleoutput" json:"simpleoutput"`
	Display      bool
}

type User struct {
	UID          string
	Username     string `bson:"username" json:"username"`
	Password     string `bson:"password" json:"password"`
	Nickname     string `bson:"nickname" json:"nickname"`
	IsChallenger bool
	Score        float64
	Privilege    string `bson:"privilege" json:"privilege"`
}

type Contest struct {
	CID string
	// pid1;pid2;pid3
	ContestName string `bson:"contestname" json:contestname`
	ProblemList string `bson:"problemlist" json:"problemlist"`
	StartTime   int64  `bson:"starttime" json:"starttime"`
	EndTime     int64  `bson:"endtime" json:"endtime"`
}

type ContestProblem struct {
	CID    string
	PID    string
	Solved int64 `bson:"solved" json:"solved"`
	Score  int64 `bson:"score" json:"score"`
}
