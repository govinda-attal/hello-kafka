.PHONY: init install-deps test build build-only pack deploy serve serve-stg ship run clean fab-gen-crypto fab-teardown fab-start

include .env
export $(shell sed 's/=.*//' .env)

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

export CUR_DIR=$(shell pwd)
export FAB_DOCKER_ENV=$(CUR_DIR)/test/dockerenv
export STG_ENV=$(CUR_DIR)/test/stgenv
export CONFIG=$(FAB_DOCKER_ENV)
export STG_CONFIG=$(STG_ENV)
export MAKE_TARGET=$@
OS?=$(shell uname -s | tr '[:upper:][:lower:]' '[:lower:][:lower:]')

init:
	go get -u github.com/golang/dep/cmd/dep

install-deps:
	rm -rf ./vendor
	dep ensure -v

test: install-deps
	ginkgo -r

build-only: 
ifndef APP_NAME
	$(error APP_NAME is not set)
endif 
	rm -rf ./dist
	mkdir dist
	mkdir dist/$(APP_NAME)
	GOOS=$(OS) GOARCH=amd64 go build -ldflags '-X "main.version=$(VERSION)"' -o ./dist/$(APP_NAME)/$(APP_NAME) ./cmd/$(APP_NAME)/...


build: test build-only
ifndef APP_NAME
	$(error APP_NAME is not set)
endif 
	
fab-gen-crypto:
	cd $(FAB_DOCKER_ENV) && ./generate.sh

fab-start:
	cd $(FAB_DOCKER_ENV) && ./start.sh

fab-teardown:
	cd $(FAB_DOCKER_ENV) && ./teardown.sh

serve:
ifndef APP_NAME
	$(error APP_NAME is not set)
endif
	cd dist/$(APP_NAME)/ && ./$(APP_NAME)

clean:
	rm -rf ./dist/ 

pack:
	docker build --build-arg X_LDFLAGS_ARGS="-X 'main.version=$(VERSION)'" -f ./build/$(APP_NAME).Dockerfile -t gattal/$(APP_NAME):$(TAG) .
	docker tag gattal/$(APP_NAME):$(TAG) gattal/$(APP_NAME):latest

upload:
	docker push gattal/$(APP_NAME):$(TAG)
	docker push gattal/$(APP_NAME):latest	

run:
	docker run --name $(APP_NAME) -p $(HOST_PORT):8080 --network=cp-all-in-one_default gattal/$(APP_NAME):$(TAG)

ship: init test pack upload clean

