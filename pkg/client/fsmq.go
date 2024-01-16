package client

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetDisposableToken(url, token) (disposableToken string, err error) {
	data := url.Values{
		"token": {token},
	}
	resp, err := http.PostForm(url+"token/get", data)
	if err != nil {
		log.Println("Failed to request token: " + err.Error())
	}

	b, err := io.ReadAll(resp.Body)
	disposableToken = strings.TrimSuffix(string(b), "\n")
	time.Sleep(200 * time.Millisecond)

	return
}

func GetBatch(url, token, queue string) (response string, err error) {
	disposableToken, err := GetDisposableToken(url, token)
	if err != nil {
		log.Println("Failed to get FSMQ token for GetBatch operation: " + err.Error())
		return
	}
	data := url.Values{
		"token": {disposableToken},
		"queue": {queue},
	}

	resp, err := http.PostForm(url+"queue/get-batch", data)
	if err != nil {
		log.Println("Failed to request batch: " + err.Error())
	}

	b, err := io.ReadAll(resp.Body)
	response = strings.TrimSuffix(string(b), "\n")
	time.Sleep(200 * time.Millisecond)

	return
}

func FsmqGetJobPayload(url, token, queue, job string) (response string, err error) {
	disposableToken, err := GetDisposableToken(url, token)
	if err != nil {
		log.Println("Failed to get FSMQ token for GetJobPayload operation: " + err.Error())
		return
	}
	data := url.Values{
		"token": {disposableToken},
		"queue": {queue},
		"job":   {job},
	}

	resp, err := http.PostForm(url+"queue/get-job-payload", data)
	if err != nil {
		log.Println("Failed to fetch job payload: " + err.Error())
	}

	b, err := io.ReadAll(resp.Body)
	response = strings.TrimSuffix(string(b), "\n")
	time.Sleep(200 * time.Millisecond)

	return
}

func FsmqAckJob(url, token, queue, job string) (err error) {
	disposableToken, err := GetDisposableToken(url, token)
	if err != nil {
		log.Println("Failed to get FSMQ token for AckJob operation: " + err.Error())
		return
	}
	data := url.Values{
		"token": {disposableToken},
		"queue": {queue},
		"job":   {job},
	}

	_, err := http.PostForm(url+"queue/ack-job", data)
	if err != nil {
		log.Println("Failed to fetch job payload: " + err.Error())
	}
	time.Sleep(200 * time.Millisecond)

	return
}

func PutJob(url, token, queue, payload string) (response string, err error) {
	disposableToken, err := GetDisposableToken(url, token)
	if err != nil {
		log.Println("Failed to get FSMQ token for PutJob operation: " + err.Error())
		return
	}
	data := url.Values{
		"token":   {disposableToken},
		"queue":   {queue},
		"payload": {payload},
	}

	resp, err := http.PostForm(url+"queue/put-job", data)
	if err != nil {
		log.Println("Failed to put job payload: " + err.Error())
		return
	}

	b, err := io.ReadAll(resp.Body)
	response = strings.TrimSuffix(string(b), "\n")

	return
}
