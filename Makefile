.PHONY: get_deps fmt
.DEFAULT_GOAL := build
build: get_deps build_server build_cli
tests: lint test
all: build build_structure build_archive docker_image
deploy: docker_image_upload

EXEC=server
CLI=cli
ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
TMP_DIR = ${ROOT}/tmp/${EXEC}
ARCHIVE=smart-home-${EXEC}.tar.gz

PROJECT ?=github.com/e154/smart-home
TRAVIS_BUILD_NUMBER ?= local
HOME ?= ${ROOT}

REV_VALUE=$(shell git rev-parse HEAD 2> /dev/null || echo "???")
REV_URL_VALUE=https://${PROJECT}/commit/${REV_VALUE}
GENERATED_VALUE=$(shell date -u +'%Y-%m-%dT%H:%M:%S%z')
DEVELOPERS_VALUE=delta54<support@e154.ru>
BUILD_NUMBER_VALUE=$(shell echo ${TRAVIS_BUILD_NUMBER})

IMAGE=smart-home-${EXEC}
DOCKER_ACCOUNT=e154
DOCKER_IMAGE_VER=${DOCKER_ACCOUNT}/${IMAGE}:${RELEASE_VERSION}
DOCKER_IMAGE_LATEST=${DOCKER_ACCOUNT}/${IMAGE}:latest

VERSION_VAR=${PROJECT}/version.VersionString
REV_VAR=${PROJECT}/version.RevisionString
REV_URL_VAR=${PROJECT}/version.RevisionURLString
GENERATED_VAR=${PROJECT}/version.GeneratedString
DEVELOPERS_VAR=${PROJECT}/version.DevelopersString
BUILD_NUMBER_VAR=${PROJECT}/version.BuildNumString
DOCKER_IMAGE_VAR=${PROJECT}/version.DockerImageString
GO_BUILD_LDFLAGS= -X ${VERSION_VAR}=${RELEASE_VERSION} -X ${REV_VAR}=${REV_VALUE} -X ${REV_URL_VAR}=${REV_URL_VALUE} -X ${GENERATED_VAR}=${GENERATED_VALUE} -X ${DEVELOPERS_VAR}=${DEVELOPERS_VALUE} -X ${BUILD_NUMBER_VAR}=${BUILD_NUMBER_VALUE} -X ${DOCKER_IMAGE_VAR}=${DOCKER_IMAGE_VER}
GO_BUILD_FLAGS= -a -installsuffix cgo -v --ldflags '${GO_BUILD_LDFLAGS}'
GO_BUILD_ENV= CGO_ENABLED=0
GO_BUILD_TAGS= -tags 'production'

test_system:
	@echo MARK: system tests
	cp ${ROOT}/conf/config.dev.json ${ROOT}/conf/config.json
	go test -v ./tests/api
	go test -v ./tests/models
	go test -v ./tests/plugins
	go test -v ./tests/scripts
	go test -v ./tests/system

test:
	@echo MARK: unit tests
	go test $(go list ./... | grep -v /tests/)
	go test -race $(go list ./... | grep -v /tests/)

install_linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.42.1

lint-todo:
	@echo MARK: make lint todo

lint:
	golangci-lint run ./...

get_deps:
	go mod tidy

fmt:
	@gofmt -l -w -s .
	@goimports -w .

cmd:
	@echo MARK: update comments
	@gocmt -i -d .

svgo:
	DIR=${ROOT}/data/icons/*
	cd ${ROOT} && svgo ${DIR} --enable=inlineStyles  --config '{ "plugins": [ { "inlineStyles": { "onlyMatchedOnce": false } }] }' --pretty

build_server:
	@echo MARK: build server
	${GO_BUILD_ENV} GOOS=linux GOARCH=amd64 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-linux-amd64
	${GO_BUILD_ENV} GOOS=linux GOARCH=arm GOARM=7 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-linux-arm-7
	${GO_BUILD_ENV} GOOS=linux GOARCH=arm GOARM=6 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-linux-arm-6
	${GO_BUILD_ENV} GOOS=linux GOARCH=arm GOARM=5 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-linux-arm-5
	${GO_BUILD_ENV} GOOS=darwin GOARCH=amd64 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-darwin-10.6-amd64

build_cli:
	@echo MARK: build cli
	cd ${ROOT}/cmd/cli && ${GO_BUILD_ENV} GOOS=linux GOARCH=amd64 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${CLI}-linux-amd64
	cd ${ROOT}/cmd/cli && ${GO_BUILD_ENV} GOOS=darwin GOARCH=amd64 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${CLI}-darwin-10.6-amd64

server:
	@echo "Building http server"
	cd ${ROOT}/api/protos/ && \
	mkdir -p ${ROOT}/api/stub && \
	protoc -I/usr/local/include -I. \
      -I${GOPATH}/src \
      -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
      -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 \
      --grpc-gateway_out=logtostderr=true:${ROOT}/api/stub \
      *.proto

	@echo "Building grpc server"
	cd ${ROOT}/api/protos/ && \
	mkdir -p ${ROOT}/api/stub && \
	protoc -I/usr/local/include -I. \
      -I${GOPATH}/src \
      -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
      -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 \
      --go-grpc_out=require_unimplemented_servers=false:${ROOT}/api/stub \
      *.proto

	@echo "Building protobuf files"
	cd ${ROOT}/api/protos/ && \
	mkdir -p ${ROOT}/api/stub && \
	protoc -I/usr/local/include -I. \
      -I${GOPATH}/src \
      -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
      -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 \
      --go_out=${ROOT}/api/stub \
      *.proto

	@echo "Building swagger.json"
	cd ${ROOT}/api/protos/ && \
	protoc -I/usr/local/include -I. \
	  -I${GOPATH}/src \
	  -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
	  -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.5.0/protoc-gen-openapiv2 \
	  -I${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0 \
	  --openapiv2_out=allow_merge=true,merge_file_name=api,logtostderr=true:${ROOT}/api \
	  *.proto

build_structure:
	@echo MARK: create app structure
	mkdir -p ${TMP_DIR}
	cd ${TMP_DIR}
	cp -r ${ROOT}/conf ${TMP_DIR}
	cp -r ${ROOT}/data ${TMP_DIR}
	cp -r ${ROOT}/snapshots ${TMP_DIR}
	cp ${ROOT}/LICENSE ${TMP_DIR}
	cp ${ROOT}/README* ${TMP_DIR}
	cp ${ROOT}/contributors.txt ${TMP_DIR}
	cp ${ROOT}/bin/docker/Dockerfile ${TMP_DIR}
	cp ${ROOT}/bin/server-installer.sh ${TMP_DIR}
	chmod +x ${TMP_DIR}/data/scripts/ping.sh
	cp ${ROOT}/${EXEC}-linux-amd64 ${TMP_DIR}
	cp ${ROOT}/${EXEC}-linux-arm-7 ${TMP_DIR}
	cp ${ROOT}/${EXEC}-linux-arm-6 ${TMP_DIR}
	cp ${ROOT}/${EXEC}-linux-arm-5 ${TMP_DIR}
	cp ${ROOT}/${EXEC}-darwin-10.6-amd64 ${TMP_DIR}
	cp ${ROOT}/${CLI}-darwin-10.6-amd64 ${TMP_DIR}
	cp ${ROOT}/${CLI}-linux-amd64 ${TMP_DIR}
	cp ${ROOT}/bin/server ${TMP_DIR}

build_archive:
	@echo MARK: build app archive
	cd ${TMP_DIR} && ls -l && tar -zcf ${HOME}/${ARCHIVE} .

build_docs:
	@echo MARK: build doc
	cd ${ROOT}/doc
	npm install postcss-cli
	hugo --gc --minify

docs_dev:
	cd ${ROOT}/doc
	hugo server --buildDrafts --verbose --source="${ROOT}/doc" --config="${ROOT}/doc/config.toml" --port=1377 --disableFastRender

doc_deploy:
	@echo MARK: deploy doc
	cd ${ROOT}/doc && \
	echo -e "node version.\n"  && \
	node -v  && \
	echo -e "npm version.\n"  && \
	npm -v  && \
	npm install -f  && \
	echo -e "hugo version.\n"  && \
	hugo version  && \
	hugo --gc --minify

	cd ${ROOT}/doc/public  && \
	git init  && \
	echo -e "Starting to documentation commit.\n"  && \
	git config --global user.email "support@e154.ru"  && \
	git config --global user.name "delta54"  && \
	git remote add upstream "https://${GITHUB_OAUTH_TOKEN}@github.com/e154/smart-home.git"  && \
	git fetch upstream  && \
	git reset upstream/gh-pages  && \
	rev=$(git rev-parse --short HEAD)  && \
	git add -A .  && \
	git commit -m "rebuild pages at ${rev}" && \
	git push -q upstream HEAD:gh-pages
	echo -e "Done documentation deploy.\n"

docker_image:
	cd ${TMP_DIR} && ls -ll && docker build -f ${ROOT}/bin/docker/Dockerfile -t ${DOCKER_ACCOUNT}/${IMAGE} .

docker_image_upload:
	echo ${DOCKER_PASSWORD} | docker login -u ${DOCKER_USERNAME} --password-stdin
	docker tag ${DOCKER_ACCOUNT}/${IMAGE} ${DOCKER_IMAGE_VER}
	echo -e "docker tag ${DOCKER_ACCOUNT}/${IMAGE} ${DOCKER_IMAGE_LATEST}"
	docker tag ${DOCKER_ACCOUNT}/${IMAGE} ${DOCKER_IMAGE_LATEST}
	docker push ${DOCKER_IMAGE_VER}
	docker push ${DOCKER_IMAGE_LATEST}

clean:
	@echo MARK: clean
	rm -rf ${TMP_DIR}
	rm -f ${ROOT}/${EXEC}-linux-amd64
	rm -f ${ROOT}/${EXEC}-linux-arm-7
	rm -f ${ROOT}/${EXEC}-linux-arm-6
	rm -f ${ROOT}/${EXEC}-linux-arm-5
	rm -f ${ROOT}/${EXEC}-darwin-10.6-amd64
	rm -f ${ROOT}/${CLI}-linux-amd64
	rm -f ${ROOT}/${CLI}-darwin-10.6-amd64
	rm -f ${HOME}/${ARCHIVE}
