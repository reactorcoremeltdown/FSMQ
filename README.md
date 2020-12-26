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

Any HTTP client capable of sending POST requests with form data can serve as FSMQ client.

TBD

## Building

I assume running `make` in the root directory of the project should suffice. Obvoiusly one needs to have `golang` installed on building machine.

Makefile TBD

---

## FAQ

**Q:** *Where are tests?*

**A:** I am rather a poor programmer, and I am still in the process of learning. I promise to get back to this topic once I learn the basics of testing in golang, but in the meantime, if you think it's unacceptable, feel free to submit a pull request, I'll happily accept!

**Q:** *What about load balancing and high availability?*

**A:** The average message throughput on typical applications of FSMQ is about 100 messages a day between a dozen of machines. For higher loads and better availability there's a number of other  solutions available.
