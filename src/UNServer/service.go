package main

import (
	//l4g "classified-lib/golang-lib/log"
	ctx "context"
	"daomanage"
	"encoding/json"
	"fmt"
	"github.com/icza/session"
	"log"
	"net/http"
	"net/url"
	"time"
	//"service/auth"
	"strconv"
)

const (
	InfoLoginFailed    = "Incorrect username or password"
	InfoLoginSucc      = "Let's GO! "
	InfoRegisterFailed = "The username has been used"

	InfoHack = "So..so?"
)

var (
	//PlbPerPage  = ctx.SvrCtx.SvrCfg.WebInfo.ProblemPerPage
	//StatPerPage = ctx.SvrCtx.SvrCfg.WebInfo.StatusPerPage
	//CstPerPage  = ctx.SvrCtx.SvrCfg.WebInfo.ContestPerPage
	PlbPerPage  int
	StatPerPage int
	CstPerPage  int
)

type StatusInfo struct {
	Status int    `bson:"status" json:"status"`
	Info   string `bson:"info" json:"info"`
}

func initSessionManager() {
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(session.NewInMemStore(), &session.CookieMngrOptions{AllowHTTP: true})
}

func initService() {
	initSessionManager()

	PlbPerPage = ctx.SvrCtx.SvrCfg.WebInfo.ProblemPerPage
	StatPerPage = ctx.SvrCtx.SvrCfg.WebInfo.StatusPerPage
	CstPerPage = ctx.SvrCtx.SvrCfg.WebInfo.ContestPerPage
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	var statusInfo StatusInfo

	sess := session.Get(r)
	log.Println("sess", sess)
	if sess != nil {
		statusInfo.Status = 200
		statusInfo.Info = sess.CAttr("username").(string)
		data, _ := json.Marshal(statusInfo)
		fmt.Fprintf(w, string(data))
		return
	} else {
		statusInfo.Status = 400
		statusInfo.Info = InfoLoginFailed
		data, _ := json.Marshal(statusInfo)
		fmt.Fprintf(w, string(data))
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	var statusInfo StatusInfo
	sess := session.Get(r)
	if sess != nil {
		session.Remove(sess, w)
		statusInfo.Status = 200
		statusInfo.Info = InfoLoginSucc
		data, _ := json.Marshal(statusInfo)
		fmt.Fprintf(w, string(data))
		return
	}
}

/*
func isAuthorized(query url.Values) bool {
	var ipaddr string

	// Authorize this query
	if len(query["ipaddr"]) > 0 {
		ipaddr = query["ipaddr"][0]
	} else {
		return false
	}
	authCode, authErr := ctx.SvrCtx.AuthMan.DoAuth(ipaddr)
	if authErr != nil {
		e := fmt.Errorf("Faild to auth this request for some internal error: %s", authErr)
		l4g.WarnLogger.Warn("%s", e)
		return false
	}
	if authCode != auth.AuthCode_OK {
		e := fmt.Errorf("Failed to auth this request, authCode: %d", authCode)
		l4g.WarnLogger.Warn("%s", e)
		return false
	}
	return true
}
*/

func registHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query      url.Values
		err        error
		statusInfo StatusInfo
	)
	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("Parse Error: Register", err)
		return
	}
	email, _ := query["email"]
	username, _ := query["username"]
	password, _ := query["password"]
	ischallenger, _ := query["ischallenger"]
	if daomanage.InsertRegister(email[0], username[0], password[0], ischallenger[0]) != nil {
		statusInfo.Status = 400
		statusInfo.Info = InfoRegisterFailed
		data, _ := json.Marshal(statusInfo)
		fmt.Fprintf(w, string(data))
		return
	}
	statusInfo.Status = 200
	statusInfo.Info = username[0]
	data, _ := json.Marshal(statusInfo)
	fmt.Fprintf(w, string(data))
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query      url.Values
		err        error
		userInfo   daomanage.User
		statusInfo StatusInfo
	)

	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("Parse Error", err)
		return
	}

	// Get user info
	if s, ok := query["username"]; ok {
		userInfo, err = daomanage.GetUserInfo("username", s[0])
		if err == nil {
			if p, ok := query["password"]; ok {
				if p[0] == userInfo.Password {
					sess := session.NewSessionOptions(&session.SessOptions{
						CAttrs: map[string]interface{}{"username": s[0]},
					})
					session.Add(sess, w)

					tNow := time.Now()
					cookie := http.Cookie{Name: "gosessionid", Value: s[0], Expires: tNow.AddDate(1, 0, 0)}
					http.SetCookie(w, &cookie)

					statusInfo.Status = 200
					statusInfo.Info = userInfo.Username
					data, _ := json.Marshal(statusInfo)
					fmt.Fprintf(w, string(data))
					return
				}
			}
		}
	}
	statusInfo.Status = 400
	statusInfo.Info = InfoLoginFailed
	data, _ := json.Marshal(statusInfo)
	fmt.Fprintf(w, string(data))
	return
}

func problemInfoHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query       url.Values
		err         error
		problemInfo daomanage.Problem
	)
	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("Parse Error", err)
		return
	}

	// Get problem info
	if s, ok := query["pid"]; ok {
		//pid, err := strconv.ParseInt(s[0], 10, 0)
		pid, err := strconv.Atoi(s[0])
		if err != nil {
			return
		}
		problemInfo, err = daomanage.GetProblemInfo(pid)
		if err == nil {
			data, _ := json.Marshal(problemInfo)
			fmt.Fprintf(w, string(data))
			return
		}
	}
	fmt.Fprintf(w, InfoHack)
}

func problemsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query url.Values
		err   error
	)
	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("Parse Error", err)
		return
	}

	var page int
	if s, ok := query["page"]; ok {
		if s[0] == "" {
			page = 1
		} else {
			//page, err = strconv.ParseInt(s[0], 10, 64)
			page, err = strconv.Atoi(s[0])
			if err != nil {
				log.Println("re")
				return
			}
		}
	}
	startIndex := PlbPerPage * (page - 1)
	endIndex := startIndex + PlbPerPage
	log.Println(PlbPerPage, startIndex, endIndex)

	// Get Problems
	problemInfoList, err := daomanage.GetProblemInRange(startIndex, endIndex)
	if err == nil {
		data, _ := json.Marshal(problemInfoList)
		fmt.Fprintf(w, string(data))
		return
	}
}

func statusInfoHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query url.Values
		err   error
	)
	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("Parse Error", err)
		return
	}

	var page int
	if s, ok := query["page"]; ok {
		if s[0] == "" {
			page = 1
		} else {
			//page, err = strconv.ParseInt(s[0], 10, 64)
			page, err = strconv.Atoi(s[0])
			if err != nil {
				log.Println("re")
				return
			}
		}
	}

	// Get status info
	startIndex := StatPerPage * (page - 1)
	endIndex := startIndex + StatPerPage
	statInfo, err := daomanage.GetStatusInRange(startIndex, endIndex)
	if err == nil {
		data, _ := json.Marshal(statInfo)
		fmt.Fprintf(w, string(data))
		return
	}
	fmt.Fprintf(w, InfoHack)
}

func contestsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query       url.Values
		err         error
		contestInfo daomanage.ContestInfoSet
	)
	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("Parse Error", err)
		return
	}

	var page int
	if s, ok := query["page"]; ok {
		if s[0] == "" {
			page = 1
		} else {
			//page, err = strconv.ParseInt(s[0], 10, 64)
			page, err = strconv.Atoi(s[0])
			if err != nil {
				log.Println("re")
				return
			}
		}
	}

	// Get contests info
	startIndex := CstPerPage * (page - 1)
	endIndex := startIndex + CstPerPage
	contestInfo, err = daomanage.GetContestInRange(startIndex, endIndex)
	if err == nil {
		data, _ := json.Marshal(contestInfo)
		fmt.Fprintf(w, string(data))
		return
	}
	fmt.Fprintf(w, InfoHack)
}

func contestInfoHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query       url.Values
		err         error
		contestInfo []daomanage.ContestProblem
	)

	type ContestProblem struct {
		CID         int    `bson:"cid" json:"cid"`
		PID         int    `bson:"pid" json:"pid"`
		ProblemName string `bson:"problemname" json:"problemname"`
		Solved      int    `bson:"solved" json:"solved"`
		Score       int    `bson:"score" json:"score"`
	}
	var contestProblemList []ContestProblem

	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("Parse Error", err)
		return
	}

	// Get contest info
	if s, ok := query["cid"]; ok {
		//cid, err := strconv.ParseInt(s[0], 10, 64)
		cid, err := strconv.Atoi(s[0])
		if err != nil {
			return
		}
		contestInfo, err = daomanage.GetContestProblems(cid)
		if err == nil {
			// data, _ := json.Marshal(contestInfo)
			for _, problem := range contestInfo {
				problemInfo, err := daomanage.GetProblemInfo(problem.PID)
				if err != nil {
					continue
				}
				singleProblemInfo := ContestProblem{problem.CID, problem.PID, problemInfo.Title, problem.Solved, problem.Score}
				contestProblemList = append(contestProblemList, singleProblemInfo)
			}
			data, _ := json.Marshal(contestProblemList)
			fmt.Printf("%s\n", data)
			fmt.Fprintf(w, string(data))
			return
		}
	}
	fmt.Fprintf(w, InfoHack)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	var statusInfo StatusInfo

	r.ParseForm()
	pid := r.PostFormValue("pid")
	code := r.PostFormValue("code")
	lang := r.PostFormValue("lang")
	author := r.PostFormValue("author")

	//npid, err := strconv.ParseInt(pid, 10, 32)
	npid, err := strconv.Atoi(pid)
	if err != nil {
		log.Println("ParseInt Error", err)
		return
	}
	daomanage.InsertSubmitQueue(npid, code, lang, author)

	statusInfo.Status = 200
	statusInfo.Info = "Submitted"
	data, _ := json.Marshal(statusInfo)
	fmt.Fprintf(w, string(data))
	return
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query    url.Values
		err      error
		userInfo daomanage.User
	)
	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("Parse Error", err)
		return
	}

	// ok := isAuthorized(query)
	// if ok == false {
	// 	fmt.Fprintf(w, InfoHack)
	// 	return
	// }

	// Get user info
	if s, ok := query["uid"]; ok {
		userInfo, err = daomanage.GetUserInfo("uid", s[0])
		if err == nil {
			data, _ := json.Marshal(userInfo)
			fmt.Fprintf(w, string(data))
			return
		}
	}
	fmt.Fprintf(w, InfoHack)
}
