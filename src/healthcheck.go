package main

import (
	"fmt"
	"net/http"
)

func Healthcheck(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "Ok, FSMQ server is running\n\nVersion: "+Version+"\nBuild date: "+BuildDate+"\nCommit ID: "+CommitID)
}
