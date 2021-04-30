package main

import (
    "fmt"
    "net/http"
    "log"
    "encoding/json"
    "io/ioutil"
    "os"
    "strconv"
    "github.com/satori/go.uuid"
)

type User struct {
    Name string
    Key string
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
        }
        authtoken := req.Form.Get("token")

        var bin []byte
        if bin, err = ioutil.ReadFile(os.Getenv("FSMQ_ACL_FILE")); err != nil {
            log.Fatalln("Can't read user database: " + err.Error())
        }

        var users []User
        err = json.Unmarshal(bin, &users)
        if err != nil {
            log.Fatalln("Can't read config file")
        }

        authSucceeded := false
        for _, element := range users {
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

                _, err = tokenFile.WriteString("{\"producer\":" + strconv.FormatBool(element.Producer) +"}")
                if err != nil {
                    log.Fatalln("Cannot write token to file:" + err.Error())
                }
                fmt.Fprint(res, tokenUUID.String() + "\n")
                break
            }
        }
        if !authSucceeded {
            res.WriteHeader(403)
            fmt.Fprint(res, "Authentication failed")
        }
    }
}
