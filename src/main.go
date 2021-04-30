package main

import (
    "net/http"
    "os"
    "log"
)

func main () {
    fsmqAppPort := os.Getenv("FSMQ_APP_PORT")
    if fsmqAppPort == "" {
        fsmqAppPort = "8080"
    }
    fsmqRoutePrefix := os.Getenv("FSMQ_ROUTE_PREFIX")

    http.HandleFunc(fsmqRoutePrefix + "/token/get", GetToken)
    http.HandleFunc(fsmqRoutePrefix + "/healthcheck", Healthcheck)
    http.HandleFunc(fsmqRoutePrefix + "/queue/get-job", QueueEndpoint)
    http.HandleFunc(fsmqRoutePrefix + "/queue/get-batch", QueueEndpoint)
    http.HandleFunc(fsmqRoutePrefix + "/queue/get-job-payload", QueueEndpoint)
    http.HandleFunc(fsmqRoutePrefix + "/queue/ack-job", QueueEndpoint)
    http.HandleFunc(fsmqRoutePrefix + "/queue/lock-job", QueueEndpoint)
    http.HandleFunc(fsmqRoutePrefix + "/queue/unlock-job", QueueEndpoint)
    http.HandleFunc(fsmqRoutePrefix + "/queue/put-job", QueueEndpoint)
    log.Fatalln(http.ListenAndServe(":" + fsmqAppPort, nil))
}
