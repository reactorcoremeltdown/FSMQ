package main

import (
    "net/http"
    "fmt"
)

func Healthcheck(res http.ResponseWriter, req *http.Request) {
    fmt.Fprint(res, "Ok, FSMQ server is running\n")
}
