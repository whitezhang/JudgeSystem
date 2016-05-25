package main

import (
	l4g "classified-lib/golang-lib/log"
	ctx "context"
	"daomanage"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"service/auth"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var query url.Values
	var ipaddr string
	var err error

	var (
		serviceName string
	// uid         string
	)

	query, err = url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println("Parse Error", err)
		return
	}

	// Authorize this query
	if len(query["ipaddr"]) > 0 {
		ipaddr = query["ipaddr"][0]
		log.Println(ipaddr)
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

	// Main section
	if s, ok := query["service"]; ok {
		serviceName = s[0]
	}

	// if s, ok := query["uid"]; ok {
	// 	uid = s[0]
	// }

	if serviceName == "getUser" {
		daomanage.GetUserInfo("tj110")
	}

	fmt.Fprintf(w, "hello")
}
