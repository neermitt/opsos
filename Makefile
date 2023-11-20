TEST?=$$(go list ./... | grep -v 'vendor')
SHELL := /bin/bash
#GOOS=darwin
GOOS=linux
GOARCH=amd64
VERSION=test

# List of targets the `readme` target should call before generating the readme
export README_DEPS ?= docs/targets.md

-include $(shell curl -sSL -o .build-harness "https://cloudposse.tools/build-harness"; echo .build-harness)

.DEFAULT_GOAL := help/short

## Lint terraform code
lint:
	$(SELF) terraform/install terraform/get-modules terraform/get-plugins terraform/lint terraform/validate

get:
	go get

# Build opsos binary
build: get
	env GOOS=${GOOS} GOARCH=${GOARCH} go build -o build/opsos -v -ldflags "-X 'github.com/neermitt/opsos/cmd.Version=${VERSION}'"

version: build
	chmod +x ./build/opsos
	./build/opsos version

deps:
	go mod download

# Run acceptance tests
testacc: get
	go test $(TEST) -v $(TESTARGS) -timeout 2m

## Terraform test
terraform/test: build
	./build/opsos terraform plan orgs/cp/tenant2/dev/us-east-2 top-level-component1
	./build/opsos terraform apply orgs/cp/tenant2/dev/us-east-2 top-level-component1 --use-plan
	./build/opsos terraform refresh orgs/cp/tenant2/dev/us-east-2 top-level-component1
	./build/opsos helmfile apply orgs/cp/tenant2/dev/us-east-2 echo-server
	./build/opsos terraform plan orgs/cp/tenant2/dev/us-east-2 top-level-component1	-- -destroy
	./build/opsos terraform apply orgs/cp/tenant2/dev/us-east-2 top-level-component1 --use-plan
	./build/opsos terraform init orgs/cp/tenant2/dev/us-east-2 top-level-component1
	./build/opsos terraform clean orgs/cp/tenant2/dev/us-east-2 top-level-component1

terraform/test2: build
	# ./build/opsos component init terraform infra/account-map
	./build/opsos stack init orgs/cp/tenant2/dev/us-east-2

terraform/test-kind: build
	./build/opsos terraform plan orgs/cp/tenant2/dev/us-east-2 infra/k8s
	./build/opsos terraform apply orgs/cp/tenant2/dev/us-east-2 infra/k8s --use-plan
	./build/opsos helmfile apply orgs/cp/tenant2/dev/us-east-2 echo-server


.PHONY: lint get build deps version testacc
