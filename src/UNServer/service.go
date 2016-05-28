package main

import (
	l4g "classified-lib/golang-lib/log"
	ctx "context"
	"daomanage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"service/auth"
)

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

func problemHandler(w http.ResponseWriter, r *http.Request) {
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

	ok := isAuthorized(query)
	if ok == false {
		fmt.Fprintf(w, "**** GO AWAY!")
	}

	// Get problem info
	if s, ok := query["pid"]; ok {
		problemInfo, err = daomanage.GetProblemInfo(s[0])
		if err == nil {
			data, _ := json.Marshal(problemInfo)
			fmt.Fprintf(w, string(data))
			return
		}
	}
	fmt.Fprintf(w, "What the fk r u looking for?")
}

func contestHandler(w http.ResponseWriter, r *http.Request) {
	var (
		query       url.Values
		err         error
		contestInfo daomanage.Contest
	)
	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("Parse Error", err)
		return
	}

	ok := isAuthorized(query)
	if ok == false {
		fmt.Fprintf(w, "**** GO AWAY!")
	}

	// Get contest info
	if s, ok := query["cid"]; ok {
		contestInfo, err = daomanage.GetContest(s[0])
		if err == nil {
			data, _ := json.Marshal(contestInfo)
			fmt.Fprintf(w, string(data))
			return
		}
	}
	fmt.Fprintf(w, "What the fk r u looking for?")
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

	ok := isAuthorized(query)
	if ok == false {
		fmt.Fprintf(w, "**** GO AWAY!")
	}

	// Get user info
	if s, ok := query["uid"]; ok {
		userInfo, err = daomanage.GetUserInfo(s[0])
		if err == nil {
			data, _ := json.Marshal(userInfo)
			fmt.Fprintf(w, string(data))
			return
		}
	}
	fmt.Fprintf(w, "What the fk r u looking for?")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var query url.Values
	var ipaddr string
	var err error

	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("Parse Error", err)
		return
	}

	// Authorize this query
	if len(query["ipaddr"]) > 0 {
		ipaddr = query["ipaddr"][0]
	} else {
		return
	}
	authCode, authErr := ctx.SvrCtx.AuthMan.DoAuth(ipaddr)
	if authErr != nil {
		e := fmt.Errorf("Faild to auth this request for some internal error: %s", authErr)
		l4g.WarnLogger.Warn("%s", e)
		return
	}
	if authCode != auth.AuthCode_OK {
		e := fmt.Errorf("Failed to auth this request, authCode: %d", authCode)
		l4g.WarnLogger.Warn("%s", e)
		return
	}

	fmt.Fprintf(w, "Hello index!")
}
