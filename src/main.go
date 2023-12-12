package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"sigs.k8s.io/yaml"
)

type fsmqWebhookHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type fsmqWebhookConfig struct {
	Name    string              `json:"name"`
	Queue   string              `json:"queue"`
	Url     string              `json:"url"`
	Headers []fsmqWebhookHeader `json:"headers"`
	Data    string              `json:"data"`
}

type fsmqAclConfig struct {
	Username string `json:"username"`
	Key      string `json:"key"`
	Producer bool   `json:"producer"`
}

type fsmqPoolConfig struct {
	Token string `json:"token"`
	Queue string `json:"queue"`
}

type fsmqNetworkConfig struct {
	Port        int    `json:"port"`
	RoutePrefix string `json:"route_prefix"`
}

type fsmqConfig struct {
	Network fsmqNetworkConfig   `json:"network"`
	Pool    fsmqPoolConfig      `json:"pool"`
	Acl     []fsmqAclConfig     `json:"acl"`
	Webhook []fsmqWebhookConfig `json:"webhook"`
}

var fsmqLogLevel string = os.Getenv("FSMQ_LOG_LEVEL")
var Version, CommitID, BuildDate string
var config fsmqConfig

func main() {
	fsmqConfigFile := os.Getenv("FSMQ_CONFIG_FILE")
	if fsmqConfigFile == "" {
		fsmqConfigFile = "/etc/fsmq/fsmq.yaml"
	}

	file, err := os.ReadFile(fsmqConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc(config.Network.RoutePrefix+"/token/get", GetToken)
	http.HandleFunc(config.Network.RoutePrefix+"/healthcheck", Healthcheck)
	http.HandleFunc(config.Network.RoutePrefix+"/queue/get-job", QueueEndpoint)
	http.HandleFunc(config.Network.RoutePrefix+"/queue/get-batch", QueueEndpoint)
	http.HandleFunc(config.Network.RoutePrefix+"/queue/get-job-payload", QueueEndpoint)
	http.HandleFunc(config.Network.RoutePrefix+"/queue/ack-job", QueueEndpoint)
	http.HandleFunc(config.Network.RoutePrefix+"/queue/lock-job", QueueEndpoint)
	http.HandleFunc(config.Network.RoutePrefix+"/queue/unlock-job", QueueEndpoint)
	http.HandleFunc(config.Network.RoutePrefix+"/queue/put-job", QueueEndpoint)
	log.Fatalln(http.ListenAndServe(":"+strconv.Itoa(config.Network.Port), nil))
}
