# FSMQ
A message queue backed by filesystem

## Motivation

A while ago I came up with an idea of an "internal" API for personal needs. Later on I found out that most of endpoints have overlapping features and do nothing other than providing the same payload to various consumers. After some brainstorming, the idea of making a message queue was born.

## USAGE

FSMQ requires some environment variables to be set to work correctly, and also an ACL provider: a JSON file that contains credentials and permissions.

One can run FSMQ both as a standalone binary, or in a container.

### Required environment variables

+ `FSMQ_ACL_FILE` — path to an ACL JSON document
+ `FSMQ_TOKEN_POOL_PATH` — path to directory where one-time tokens are stored
+ `FSMQ_QUEUE_POOL_PATH` — path to directory where queues(subdirectories) and messages are stored

### Optional environment variables

+ `FSMQ_APP_PORT` — a port number for app to bind to. Defaults to 8080 if unset.
+ `FSMQ_ROUTE_PREFIX` — a prefix for all routes. Defaults to empty string, can be useful to serve FSMQ as a sub-route of another application.

### ACL file structure

TBD

### Launching FSMQ

TBD

### Clients

Any HTTP client capable of sending POST requests with form data can serve as FSMQ client. Below you can find the list of endpoints and parameters that FSMQ can accept and process.

+ `GET` `/<prefix>/healthcheck` — Returns a string, indicates that the application is running.
+ `POST` `/<prefix>/token/get` — Returns one-time token for authorizing on other endpoints.
    Accepts following form parameters:
    + `token=` — A token from ACL file.
    Example: `curl -XPOST --data-urlencode "token=foo" https://fsmq.tld/<prefix>/token/get`
+ `POST` `/<prefix>/queue/get-job` — Returns the ID of a next job available for processing in a given queue.
    Accepts following form parameters:
    + `token=` — One-time token produced by `/<prefix>/token/get` endpoint
    + `queue=` — An ID of a queue
    Example: `curl -XPOST --data-urlencode "token=foo" --data-urlencode "queue=bar" https://fsmq.tld/<prefix>/queue/get-job`
+ `POST` `/<prefix>/queue/get-job-payload` — Returns a payload (file contents) of a specified job of a specified queue.
    Accepts following form parameters:
    + `token=` — One-time token produced by `/<prefix>/token/get` endpoint
    + `queue=` — An ID of a queue
    + `job=` — An ID of a job
    Example: `curl -XPOST --data-urlencode "token=foo" --data-urlencode "queue=bar" --data-urlencode "job=baz" https://fsmq.tld/<prefix>/queue/get-job-payload`
+ `POST` `/<prefix>/queue/lock-job` — Locks a job (renames a file by adding `-locked` suffix) to protect it from being processed by other consumers. Returns `OK` on success.
    Accepts following form parameters:
    + `token=` — One-time token produced by `/<prefix>/token/get` endpoint
    + `queue=` — An ID of a queue
    + `job=` — An ID of a job
    Example: `curl -XPOST --data-urlencode "token=foo" --data-urlencode "queue=bar" --data-urlencode "job=baz" https://fsmq.tld/<prefix>/queue/lock-job`
+ `POST` `/<prefix>/queue/unlock-job` — Makes any locked job available for processing by other consumers. Returns `OK` on success.
    Accepts following form parameters:
    + `token=` — One-time token produced by `/<prefix>/token/get` endpoint
    + `queue=` — An ID of a queue
    + `job=` — An ID of a job
    Example: `curl -XPOST --data-urlencode "token=foo" --data-urlencode "queue=bar" --data-urlencode "job=baz" https://fsmq.tld/<prefix>/queue/unlock-job`
+ `POST` `/<prefix>/queue/ack-job` — Permanently removes (acknowledges) job from the queue. Returns `OK` on success.
    Accepts following form parameters:
    + `token=` — One-time token produced by `/<prefix>/token/get` endpoint
    + `queue=` — An ID of a queue
    + `job=` — An ID of a job
    Example: `curl -XPOST --data-urlencode "token=foo" --data-urlencode "queue=bar" --data-urlencode "job=baz" https://fsmq.tld/<prefix>/queue/ack-job`
+ `POST` `/<prefix>/queue/put-job` — Produce a job & make it available to other consumers. Returns a job ID on success.
    + `token=` — One-time token produced by `/<prefix>/token/get` endpoint
    + `queue=` — An ID of a queue
    + `payload=` — A message for consumers to process
    Example: `curl -XPOST --data-urlencode "token=foo" --data-urlencode "queue=bar" --data-urlencode "payload=qux" https://fsmq.tld/<prefix>/queue/put-job`

## Building

I assume running `make` in the root directory of the project should suffice. Obvoiusly one needs to have `golang` installed on building machine.

Makefile TBD

---

## FAQ

**Q:** *Where are tests?*

**A:** I am rather a poor programmer, and I am still in the process of learning. I promise to get back to this topic once I learn the basics of testing in golang, but in the meantime, if you think it's unacceptable, feel free to submit a pull request, I'll happily accept!

**Q:** *What about load balancing and high availability?*

**A:** The average message throughput on typical applications of FSMQ is about 100 messages a day between a dozen of machines. For higher loads and better availability there's a number of other  solutions available.
