SHELL := /bin/bash

# ==============================================================================
# Testing running system

# For testing a simple query on the system. Don't forget to `make seed` first.
# curl --user "admin@example.com:gophers" http://localhost:3000/v1/users/token
# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2
#
# For testing load on the service.
# go install github.com/rakyll/hey@latest
# hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2
#
# Access metrics directly (4000) or through the sidecar (3001)
# go install github.com/divan/expvarmon@latest
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# expvarmon -ports=":3001" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
#
# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem
# ./wakt-admin genkey
#
# Testing coverage.
# go test -coverprofile p.out
# go tool cover -html p.out
#
# Test debug endpoints.
# curl http://localhost:4000/debug/liveness
# curl http://localhost:4000/debug/readiness
#
# Running pgcli client for database.
# brew install pgcli
# pgcli postgresql://postgres:postgres@localhost
#
# Launch zipkin.
# http://localhost:9411/zipkin/


# ==============================================================================
# Install dependencies

dev.setup.mac:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize

# ==============================================================================
# Building containers

# $(shell git rev-parse --short HEAD)
VERSION := 0.1.0

all: wakt metrics

wakt:
	docker build \
		-f zarf/docker/dockerfile.wakt-api \
		-t wakt-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

metrics:
	docker build \
		-f zarf/docker/dockerfile.metrics \
		-t metrics-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := shaef-starter-cluster

# Upgrade to latest Kind (>=v0.11): e.g. brew upgrade kind
# For full Kind v0.11 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.0
# Kind release used for our project: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.1
# The image used below was copied by the above link and supports both amd64 and arm64.

kind-up:
	kind create cluster \
		--image kindest/node:v1.23.5\
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=wakt-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	cd zarf/k8s/kind/wakt-pod; kustomize edit set image wakt-api-image=wakt-api-amd64:$(VERSION)
	cd zarf/k8s/kind/wakt-pod; kustomize edit set image metrics-image=metrics-amd64:$(VERSION)
	kind load docker-image wakt-api-amd64:$(VERSION) --name $(KIND_CLUSTER)
	kind load docker-image metrics-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=480s --for=condition=Available deployment/database-pod
	kustomize build zarf/k8s/kind/zipkin-pod | kubectl apply -f -
	kubectl wait --namespace=zipkin-system --timeout=480s --for=condition=Available deployment/zipkin-pod
	kustomize build zarf/k8s/kind/wakt-pod | kubectl apply -f -

kind-services-delete:
	kustomize build zarf/k8s/kind/wakt-pod | kubectl delete -f -
	kustomize build zarf/k8s/kind/zipkin-pod | kubectl delete -f -
	kustomize build zarf/k8s/kind/database-pod | kubectl delete -f -

kind-restart:
	kubectl rollout restart deployment wakt-pod

kind-update: all kind-load kind-restart

kind-all-load-apply: all kind-load kind-apply

kind-logs:
	kubectl logs -l app=wakt --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go

kind-logs-wakt:
	kubectl logs -l app=wakt --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go -service=wakt-API

kind-logs-metrics:
	kubectl logs -l app=wakt --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go -service=METRICS

kind-logs-db:
	kubectl logs -l app=database --namespace=database-system --all-containers=true -f --tail=100

kind-logs-zipkin:
	kubectl logs -l app=zipkin --namespace=zipkin-system --all-containers=true -f --tail=100

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-wakt:
	kubectl get pods -o wide --watch --namespace=wakt-system

kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system

kind-status-zipkin:
	kubectl get pods -o wide --watch --namespace=zipkin-system

kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=wakt

kind-describe-deployment:
	kubectl describe deployment wakt-pod

kind-describe-replicaset:
	kubectl get rs
	kubectl describe rs -l app=wakt

kind-events:
	kubectl get ev --sort-by metadata.creationTimestamp

kind-events-warn:
	kubectl get ev --field-selector type=Warning --sort-by metadata.creationTimestamp

kind-context-wakt:
	kubectl config set-context --current --namespace=wakt-system

kind-shell:
	kubectl exec -it $(shell kubectl get pods | grep wakt | cut -c1-26) --container wakt-api -- /bin/sh

kind-database:
	# ./admin --db-disable-tls=1 migrate
	# ./admin --db-disable-tls=1 seed

# ==============================================================================
# Administration

migrate:
	go run app/tooling/wakt-admin/main.go migrate

seed: migrate
	go run app/tooling/wakt-admin/main.go seed

genkey:
	go run app/tooling/wakt-admin/main.go genkey

# ==============================================================================
# Running tests within the local computer
run:
	go run app/services/wakt-api/main.go
runfmt:
	go run app/services/wakt-api/main.go | app/tooling/logfmt/main.go
test:
	go test ./... -count=1
	staticcheck -checks=all ./...

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Docker support

docker-down:
	docker rm -f $(shell docker ps -aq)

docker-clean:
	docker system prune -f	

docker-kind-logs:
	docker logs -f $(KIND_CLUSTER)-control-plane

