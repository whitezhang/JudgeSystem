package daomanage

import ()

type Problem struct {
	PID          int64
	Title        string `bson:"title" json:"title"`
	Description  string `bson:"description" json:"description"`
	Time         string `bson:"time" json:"time"`
	Memory       string `bson:"memory" json:"memory"`
	Input        string `bson:"input" json:"input"`
	Output       string `bson:"output" json:"output"`
	SimpleInput  string `bson:"simpleinput" json:"simpleinput"`
	SimpleOutput string `bson:"simpleoutput" json:"simpleoutput"`
	Solved       int64
	Display      bool
}

type User struct {
	UID          int64
	Username     string `bson:"username" json:"username"`
	Password     string `bson:"password" json:"password"`
	Nickname     string `bson:"nickname" json:"nickname"`
	IsChallenger bool
	Score        float64
	Privilege    string `bson:"privilege" json:"privilege"`
}

type Contest struct {
	CID int64
	// pid1;pid2;pid3
	ContestName string `bson:"contestname" json:contestname`
	ProblemList string `bson:"problemlist" json:"problemlist"`
	StartTime   int64  `bson:"starttime" json:"starttime"`
	EndTime     int64  `bson:"endtime" json:"endtime"`
}

type ContestProblem struct {
	CID    int64
	PID    int64
	Solved int64 `bson:"solved" json:"solved"`
	Score  int64 `bson:"score" json:"score"`
}

// Specific type
type ProblemInfo struct {
	PID         int64  `bson:"pid" json:"pid"`
	ProblemName string `bson:"title" json:"title"`
	Solved      int64  `bson:"solved" json:"solved"`
}
