package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	Name     string
	Key      string
	Producer bool
}

func GetToken(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(503)
		fmt.Fprint(res, "This endpoint processes POST requests only")
	} else {
		err := req.ParseForm()
		if err != nil {
			log.Println("Failed to parse form data: " + err.Error())
		} else {
			if fsmqLogLevel == "debug" {
				for key, value := range req.Form {
					log.Printf("%s = %s\n", key, value)
				}
			}
		}
		authtoken := req.Form.Get("token")

		authSucceeded := false
		for _, element := range config.Acl {
			if element.Key == authtoken {
				authSucceeded = true
				tokenUUID := uuid.NewV4()
				if err != nil {
					log.Fatalln("Cannot issue UUIDv4 token: " + err.Error())
				}
				tokenFile, err := os.Create(os.Getenv("FSMQ_TOKEN_POOL_PATH") + "/token-" + tokenUUID.String())
				if err != nil {
					log.Fatalln("Cannot create token file: " + err.Error())
				}

				defer tokenFile.Close()

				_, err = tokenFile.WriteString("{\"producer\":" + strconv.FormatBool(element.Producer) + "}")
				if err != nil {
					log.Fatalln("Cannot write token to file:" + err.Error())
				}
				fmt.Fprint(res, tokenUUID.String()+"\n")
				break
			}
		}
		if !authSucceeded {
			res.WriteHeader(403)
			fmt.Fprint(res, "Authentication failed")
		}
	}
}
