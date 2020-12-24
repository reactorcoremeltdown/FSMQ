package main

import (
    "os"
    "log"
    "encoding/json"
    "io/ioutil"
)

type authRole struct {
    Producer bool
}

func CheckAuth (token string) (success, producer bool) {
    filepath := os.Getenv("FSMQ_TOKEN_POOL_PATH") + "/token-" + token
    b, err := ioutil.ReadFile(filepath)
    if err != nil {
        log.Println("Cannot read file: " + err.Error())
        success = false
        producer = false
    }
    var role authRole
    err = json.Unmarshal(b, &role)
    if err != nil {
        log.Println("Cannot read role from token file:" + err.Error())
        success = false
        producer = false
    }
    err = os.Remove(filepath)
    if err != nil {
        log.Println("Cannot remove token file: " + err.Error())
        success = false
        producer = false
    } else {
        success = true
        producer = role.Producer
    }

    return
}
