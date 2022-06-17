# Wakt
![](wakt.png)

----
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/AhmedShaef/wakt)](https://github.com/AhmedShaef/wakt)
[![Wiki](https://img.shields.io/badge/wiki-wakt-blue.svg)](https://github.com/AhmedShaef/wakt/wiki)
![License](https://img.shields.io/badge/license-GNU%3A%20General%20Public%20License-blue.svg)

# Getting Started
Wakt is an open-source time tracking Microservice, based on [Ardanlabs service 3](https://github.com/ardanlabs/service) and inspired by [Toggl track](https://toggl.com/track/).

## fast run (this just run wakt api)
1- the Go 1.18 + and postgres 14.2 is requred.

2- use commend to run main.go.
```shell
make run
```
3- apply the scemma and seed demo data.
```shell
make seed
```
### Prerequisites

* [Go 1.18 +](https://golang.org/doc/install)
* [docker](https://www.docker.com/community-edition)
* [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/)
* [kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)

You can run this make command to use brew to install all the software above.
```shell
make dev.setup.mac
```
## To start using wakt
```shell
    make all
    make kind-up
    make kind-load
    make kind-apply  
```   
## Check services status  

### Check status
```shell
    make kind-status 
```
### Check logs
```shell
    make kind-logs
```
also you can log specific service ex.
```shell
    make kind-logs-wakt
```
### Check traces

Use Zipkin to query traces in [localhost:9411](http://localhost:9411)

### Check metrics

use the expvar sidecar service in port 4000

Ckeck readiness in [localhost:4000/debug/readiness](http://localhost:4000/debug/readiness) the status ok mean the api up and running

Liveness is in [locakhost:4000/debug/liveness](http://locakhost/400/debug/liveness) if you get data that's mean the db is running and connected
    
## To stop using wakt
```shell
    make kind-down
```
## Run tests (that support unit and integrated test)
```shell
    make test
```
----
## Data Model
![](data-model.png)
----
## TODO for v1.0.0
### in API :
1. [ ] Decrase number of db cooniction per request
2. [ ] Improve notufcation system
3. [ ] Report package
4. [ ] invoce and payment
5. [ ] Oauth2
6. [ ] improve comments
7. [ ] missing funcalities

### in UI:
* Design UI/UX and Code (using React.js) the following:
  1. [ ] App Dashboard
  2. [ ] Home Page
  3. [ ] Price Page
  4. [ ] SignUp/Login/forget password Pages
  5. [ ] Email HTML Template for:
     1. [ ] SignUp validation
     2. [ ] Invitations
     3. [ ] Reset email/password
     4. [ ] Reports and invoices
