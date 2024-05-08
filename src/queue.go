package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/csotherden/strftime"
	uuid "github.com/satori/go.uuid"
)

func sortFiles(jobs []os.FileInfo) {
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].ModTime().Unix() < jobs[j].ModTime().Unix()
	})
}

func queueGetJob(queueID string) (string, error) {
	jobs, err := ioutil.ReadDir(config.Pool.Queue + "/" + queueID + "/")
	if err != nil {
		log.Println("Cannot access token pool:" + err.Error())
		return "", fmt.Errorf("EOQ")
	}

	if len(jobs) != 0 {
		sortFiles(jobs)
		for _, job := range jobs {
			if !strings.Contains(job.Name(), "-locked") {
				return job.Name(), nil
			}
		}
		return "", fmt.Errorf("EOQ")
	} else {
		return "", fmt.Errorf("EOQ")
	}
}

func queueGetBatch(queueID string) (string, error) {
	jobs, err := ioutil.ReadDir(config.Pool.Queue + "/" + queueID + "/")
	if err != nil {
		log.Println("Cannot access token pool:" + err.Error())
		return "", fmt.Errorf("EOQ")
	}

	output := ""
	if len(jobs) != 0 {
		sortFiles(jobs)
		for _, job := range jobs {
			if !strings.Contains(job.Name(), "-locked") {
				output = output + job.Name() + "\n"
			}
		}
		if output != "" {
			return output, nil
		} else {
			return "", fmt.Errorf("EOQ")
		}
	} else {
		return "", fmt.Errorf("EOQ")
	}
}

func queueGetJobPayload(queueID, jobID string) (string, error) {
	bin, err := ioutil.ReadFile(config.Pool.Queue + "/" + queueID + "/" + jobID)
	if err != nil {
		log.Println("Cannot read job payload: " + err.Error())
		return "", err
	}

	return string(bin), nil
}

func queueAckJob(queueID, jobID string) error {
	err := os.Remove(config.Pool.Queue + "/" + queueID + "/" + jobID)
	if err != nil {
		log.Println("Cannot remove job: " + err.Error())
		return err
	}
	return nil
}

func queueLockJob(queueID, jobID string) error {
	err := os.Rename(config.Pool.Queue+"/"+queueID+"/"+jobID,
		config.Pool.Queue+"/"+queueID+"/"+jobID+"-locked")
	if err != nil {
		log.Println("Cannot lock job: " + err.Error())
		return err
	}
	return nil
}

func queueUnlockJob(queueID, jobID string) error {
	err := os.Rename(config.Pool.Queue+"/"+queueID+"/"+jobID+"-locked",
		config.Pool.Queue+"/"+queueID+"/"+jobID)
	if err != nil {
		log.Println("Cannot unlock job: " + err.Error())
		return err
	}
	return nil
}

func queuePutJob(queueID, jobPayload string) (string, error) {
	jobUUID := uuid.NewV4()
	if _, err := os.Stat(config.Pool.Queue + "/" + queueID); os.IsNotExist(err) {
		err = os.Mkdir(config.Pool.Queue+"/"+queueID, 0755)
		if err != nil {
			log.Println("Cannot create queue: " + err.Error())
			return "", err
		}
	}
	jobFilename := strftime.Format("%Y-%m-%d_%H-%M-%S_", time.Now()) +
		jobUUID.String()
	jobFilepath := config.Pool.Queue +
		"/" +
		queueID +
		"/" +
		jobFilename
	jobFile, err := os.Create(jobFilepath)
	if err != nil {
		log.Println("Cannot create job file: " + err.Error())
		return "", err
	}

	defer jobFile.Close()

	_, err = jobFile.WriteString(jobPayload)
	if err != nil {
		log.Println("Cannot write payload to file:" + err.Error())
		return "", err
	}

	for _, webhook := range config.Webhook {
		if webhook.Queue == queueID {
			log.Println("Firing " + webhook.Name + " webhook")
			req, err := http.NewRequest("POST", webhook.Url, bytes.NewBuffer([]byte(webhook.Data)))
			for _, header := range webhook.Headers {
				req.Header.Set(header.Key, header.Value)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Println("Failed to run a webhook" + err.Error())
			}
			defer resp.Body.Close()

			log.Println("Webhook "+webhook.Name+" response Status:", resp.Status)
			log.Println("Webhook "+webhook.Name+" response Headers:", resp.Header)
			body, _ := io.ReadAll(resp.Body)
			log.Println("Webhook "+webhook.Name+" response Body:", string(body))
		}
	}
	return jobFilename, nil
}

func QueueEndpoint(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(503)
		fmt.Fprint(res, "This endpoint processes POST requests only\n")
	} else {
		err := req.ParseForm()
		if err != nil {
			log.Println("Failed to parse form data: " + err.Error())
		}
		authtoken := req.Form.Get("token")

		authSucceeded, producer := CheckAuth(authtoken)

		if authSucceeded {
			queueID := req.Form.Get("queue")
			if queueID == "" {
				res.WriteHeader(400)
				fmt.Fprint(res, "This endpoint needs a queue ID specified\n")
			} else {
				switch path := req.URL.Path; path {
				case config.Network.RoutePrefix + "/queue/get-job":
					jobID, err := queueGetJob(queueID)
					if err != nil {
						res.WriteHeader(404)
						fmt.Fprint(res, err.Error())
					} else {
						fmt.Fprint(res, jobID)
					}
				case config.Network.RoutePrefix + "/queue/get-batch":
					batch, err := queueGetBatch(queueID)
					if err != nil {
						res.WriteHeader(404)
						fmt.Fprint(res, err.Error())
					} else {
						fmt.Fprint(res, batch)
					}
				case config.Network.RoutePrefix + "/queue/get-job-payload":
					jobID := req.Form.Get("job")
					if jobID == "" {
						res.WriteHeader(400)
						fmt.Fprint(res, "This endpoint needs a job ID specified\n")
					} else {
						jobPayload, err := queueGetJobPayload(queueID, jobID)
						if err != nil {
							res.WriteHeader(404)
							fmt.Fprint(res, "No such job\n")
						} else {
							fmt.Fprint(res, jobPayload)
						}
					}
				case config.Network.RoutePrefix + "/queue/ack-job":
					jobID := req.Form.Get("job")
					if jobID == "" {
						res.WriteHeader(400)
						fmt.Fprint(res, "This endpoint needs a job ID specified\n")
					} else {
						err = queueAckJob(queueID, jobID)
						if err != nil {
							res.WriteHeader(503)
							fmt.Fprint(res, "Failed to acknowledge job "+jobID+"\n")
						} else {
							fmt.Fprint(res, "OK\n")
						}
					}
				case config.Network.RoutePrefix + "/queue/lock-job":
					jobID := req.Form.Get("job")
					if jobID == "" {
						res.WriteHeader(400)
						fmt.Fprint(res, "This endpoint needs a job ID specified\n")
					} else {
						err = queueLockJob(queueID, jobID)
						if err != nil {
							res.WriteHeader(503)
							fmt.Fprint(res, "Failed to lock job "+jobID+"\n")
						} else {
							fmt.Fprint(res, "OK\n")
						}
					}
				case config.Network.RoutePrefix + "/queue/unlock-job":
					jobID := req.Form.Get("job")
					if jobID == "" {
						res.WriteHeader(400)
						fmt.Fprint(res, "This endpoint needs a job ID specified\n")
					} else {
						err = queueUnlockJob(queueID, jobID)
						if err != nil {
							res.WriteHeader(503)
							fmt.Fprint(res, "Failed to unlock job "+jobID+"\n")
						} else {
							fmt.Fprint(res, "OK\n")
						}
					}
				case config.Network.RoutePrefix + "/queue/put-job":
					if producer {
						jobPayload := req.Form.Get("payload")
						if jobPayload == "" {
							res.WriteHeader(400)
							fmt.Fprint(res, "Payload cannot be empty\n")
						} else {
							jobID, err := queuePutJob(queueID, jobPayload)
							if err != nil {
								res.WriteHeader(503)
								fmt.Fprint(res, "Failed to register a new job\n")
							} else {
								fmt.Fprint(res, jobID)
							}
						}
					} else {
						res.WriteHeader(403)
						fmt.Fprint(res, "You are not authorized to produce jobs\n")
					}
				default:
					res.WriteHeader(503)
					fmt.Fprint(res, "Cannot find path: "+path+"\n")
				}
			}
		} else {
			res.WriteHeader(403)
			fmt.Fprint(res, "Authentication failed\n")
		}
	}
}
