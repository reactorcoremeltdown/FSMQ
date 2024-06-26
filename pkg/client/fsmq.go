package client

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetDisposableToken(queueURL, token string) (disposableToken string, err error) {
	data := url.Values{
		"token": {token},
	}
	resp, err := http.PostForm(queueURL+"token/get", data)
	if err != nil {
		log.Println("Failed to request token: " + err.Error())
	}

	b, err := io.ReadAll(resp.Body)
	disposableToken = strings.TrimSuffix(string(b), "\n")
	time.Sleep(200 * time.Millisecond)

	return
}

func GetBatch(queueURL, token, queue string) (response string, err error) {
	disposableToken, err := GetDisposableToken(queueURL, token)
	if err != nil {
		log.Println("Failed to get FSMQ token for GetBatch operation: " + err.Error())
		return
	}
	data := url.Values{
		"token": {disposableToken},
		"queue": {queue},
	}

	resp, err := http.PostForm(queueURL+"queue/get-batch", data)
	if err != nil {
		log.Println("Failed to request batch: " + err.Error())
	}

	b, err := io.ReadAll(resp.Body)
	response = strings.TrimSuffix(string(b), "\n")
	time.Sleep(200 * time.Millisecond)

	return
}

func GetJobPayload(queueURL, token, queue, job string) (response string, err error) {
	disposableToken, err := GetDisposableToken(queueURL, token)
	if err != nil {
		log.Println("Failed to get FSMQ token for GetJobPayload operation: " + err.Error())
		return
	}
	data := url.Values{
		"token": {disposableToken},
		"queue": {queue},
		"job":   {job},
	}

	resp, err := http.PostForm(queueURL+"queue/get-job-payload", data)
	if err != nil {
		log.Println("Failed to fetch job payload: " + err.Error())
	}

	b, err := io.ReadAll(resp.Body)
	response = strings.TrimSuffix(string(b), "\n")
	time.Sleep(200 * time.Millisecond)

	return
}

func AckJob(queueURL, token, queue, job string) (err error) {
	disposableToken, err := GetDisposableToken(queueURL, token)
	if err != nil {
		log.Println("Failed to get FSMQ token for AckJob operation: " + err.Error())
		return
	}
	data := url.Values{
		"token": {disposableToken},
		"queue": {queue},
		"job":   {job},
	}

	_, err = http.PostForm(queueURL+"queue/ack-job", data)
	if err != nil {
		log.Println("Failed to ack job: " + err.Error())
	}
	time.Sleep(200 * time.Millisecond)

	return
}

func DiscardAllJobs(queueURL, token, queue string) (err error) {
	disposableToken, err := GetDisposableToken(queueURL, token)
	if err != nil {
		log.Println("Failed to get FSMQ token for AckJob operation: " + err.Error())
		return
	}
	data := url.Values{
		"token": {disposableToken},
		"queue": {queue},
	}

	_, err = http.PostForm(queueURL+"queue/discard-all-jobs", data)
	if err != nil {
		log.Println("Failed to discard all jobs in a queue: " + err.Error())
	}
	time.Sleep(200 * time.Millisecond)

	return
}

func LockJob(queueURL, token, queue, job string) (err error) {
	disposableToken, err := GetDisposableToken(queueURL, token)
	if err != nil {
		log.Println("Failed to get FSMQ token for LockJob operation: " + err.Error())
		return
	}
	data := url.Values{
		"token": {disposableToken},
		"queue": {queue},
		"job":   {job},
	}

	_, err = http.PostForm(queueURL+"queue/lock-job", data)
	if err != nil {
		log.Println("Failed to lock job: " + err.Error())
	}
	time.Sleep(200 * time.Millisecond)

	return
}

func UnlockJob(queueURL, token, queue, job string) (err error) {
	disposableToken, err := GetDisposableToken(queueURL, token)
	if err != nil {
		log.Println("Failed to get FSMQ token for LockJob operation: " + err.Error())
		return
	}
	data := url.Values{
		"token": {disposableToken},
		"queue": {queue},
		"job":   {job},
	}

	_, err = http.PostForm(queueURL+"queue/unlock-job", data)
	if err != nil {
		log.Println("Failed to unlock job: " + err.Error())
	}
	time.Sleep(200 * time.Millisecond)

	return
}

func PutJob(queueURL, token, queue, payload string) (response string, err error) {
	disposableToken, err := GetDisposableToken(queueURL, token)
	if err != nil {
		log.Println("Failed to get FSMQ token for PutJob operation: " + err.Error())
		return
	}
	data := url.Values{
		"token":   {disposableToken},
		"queue":   {queue},
		"payload": {payload},
	}

	resp, err := http.PostForm(queueURL+"queue/put-job", data)
	if err != nil {
		log.Println("Failed to put job payload: " + err.Error())
		return
	}

	b, err := io.ReadAll(resp.Body)
	response = strings.TrimSuffix(string(b), "\n")

	return
}
