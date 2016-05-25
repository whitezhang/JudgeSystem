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
	ProblemList string `bson:"problemlist" json:"problemlist"`
}
