.PHONY: get_deps fmt
.DEFAULT_GOAL := build
tests: lint test

EXEC=server
CLI=cli
ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
SERVER_DIR = ${ROOT}/tmp/${EXEC}
COMMON_DIR = ${ROOT}/tmp/common

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
RELEASE_VERSION ?= v0.0.0
DOCKER_IMAGE_VER=${DOCKER_ACCOUNT}/${IMAGE}:${RELEASE_VERSION}

VERSION_VAR=${PROJECT}/version.VersionString
REV_VAR=${PROJECT}/version.RevisionString
REV_URL_VAR=${PROJECT}/version.RevisionURLString
GENERATED_VAR=${PROJECT}/version.GeneratedString
DEVELOPERS_VAR=${PROJECT}/version.DevelopersString
BUILD_NUMBER_VAR=${PROJECT}/version.BuildNumString
DOCKER_IMAGE_VAR=${PROJECT}/version.DockerImageString

GO_BUILD_LDFLAGS=-X ${VERSION_VAR}=${RELEASE_VERSION} -X ${REV_VAR}=${REV_VALUE} -X ${REV_URL_VAR}=${REV_URL_VALUE} -X ${GENERATED_VAR}=${GENERATED_VALUE} -X ${DEVELOPERS_VAR}=${DEVELOPERS_VALUE} -X ${BUILD_NUMBER_VAR}=${BUILD_NUMBER_VALUE} -X ${DOCKER_IMAGE_VAR}=${DOCKER_IMAGE_VER}
GO_TEST=test -tags test -v

test_system:
	@echo MARK: system tests
	cp ${ROOT}/conf/config.dev.json ${ROOT}/conf/config.json
	go ${GO_TEST} ./tests/api
	go ${GO_TEST} ./tests/models
	go ${GO_TEST} ./tests/scripts
	go ${GO_TEST} ./tests/system
	go ${GO_TEST} ./tests/plugins/alexa
	go ${GO_TEST} ./tests/plugins/area
	go ${GO_TEST} ./tests/plugins/cgminer
	go ${GO_TEST} ./tests/plugins/email
	go ${GO_TEST} ./tests/plugins/messagebird
	go ${GO_TEST} ./tests/plugins/modbus_rtu
	go ${GO_TEST} ./tests/plugins/modbus_tcp
	go ${GO_TEST} ./tests/plugins/moon
	go ${GO_TEST} ./tests/plugins/node
	go ${GO_TEST} ./tests/plugins/scene
	go ${GO_TEST} ./tests/plugins/sensor
	go ${GO_TEST} ./tests/plugins/sun
	go ${GO_TEST} ./tests/plugins/telegram
	go ${GO_TEST} ./tests/plugins/trigger_alexa
	go ${GO_TEST} ./tests/plugins/trigger_empty
	go ${GO_TEST} ./tests/plugins/trigger_state
	go ${GO_TEST} ./tests/plugins/trigger_system
	go ${GO_TEST} ./tests/plugins/trigger_time
	go ${GO_TEST} ./tests/plugins/twilio
	go ${GO_TEST} ./tests/plugins/weather_met
	go ${GO_TEST} ./tests/plugins/weather_owm
	go ${GO_TEST} ./tests/plugins/zigbee2mqtt

test:
	@echo MARK: unit tests
	go ${GO_TEST} $(shell go list ./... | grep -v /tmp | grep -v /tests) -timeout 60s -race -covermode=atomic -coverprofile=coverage.out

test_without_race:
	@echo MARK: unit tests
	go ${GO_TEST} $(shell go list ./... | grep -v /tmp | grep -v /tests) -timeout 60s -covermode=atomic -coverprofile=coverage.out

install_linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2

install_typedoc:
	npm install -g jsdoc

lint-todo:
	@echo MARK: make lint todo

lint:
	golangci-lint run

get_deps:
	go mod tidy

fmt:
	@gofmt -l -w -s .
	@goimports -w .

comments:
	@echo MARK: update comments
	@gocmt -i -d .

svgo:
	DIR=${ROOT}/data/icons/*
	cd ${ROOT} && svgo ${DIR} --enable=inlineStyles  --config '{ "plugins": [ { "inlineStyles": { "onlyMatchedOnce": false } }] }' --pretty

build_public:
	@echo MARK: build public
	echo -e "node version.\n"  && \
	node -v  && \
	echo -e "npm version.\n"  && \
	npm -v  && \
	npm i -g pnpm@8.15.1  && \
	echo -e "pnpm version.\n"  && \
	pnpm -v && \
	cd ${ROOT}/static_source/admin && \
	pnpm i && \
	pnpm run build:pro && \
	rm -rf ${ROOT}/build/public && \
	mkdir -p ${ROOT}/build && \
	mv ${ROOT}/static_source/admin/dist-pro ${ROOT}/build/public

server:
	@echo "Building http server"
	mkdir -p ${ROOT}/internal/api/stub && \
	oapi-codegen -generate server -package stub ${ROOT}/internal/api/api.swagger.yaml > ${ROOT}/internal/api/stub/server.go && \
	oapi-codegen -generate types -package stub ${ROOT}/internal/api/api.swagger.yaml > ${ROOT}/internal/api/stub/types.go

build_common_structure:
	@echo MARK: create common structure
	mkdir -p ${COMMON_DIR}
	mkdir -p ${COMMON_DIR}/snapshots
	cd ${COMMON_DIR}
	cp -r ${ROOT}/conf ${COMMON_DIR}
	cp -r ${ROOT}/data ${COMMON_DIR}
	cp ${ROOT}/LICENSE ${COMMON_DIR}
	cp ${ROOT}/README* ${COMMON_DIR}
	cp ${ROOT}/CONTRIBUTING.md ${COMMON_DIR}
	chmod +x ${COMMON_DIR}/data/scripts/ping.sh
	cd ${COMMON_DIR} && ls -l && tar -zcf ${ROOT}/common.tar.gz .

build_docs:
	@echo MARK: build doc
	cd ${ROOT}/doc
	npm install postcss-cli
	hugo --gc --minify

docs_dev:
	cd ${ROOT}/doc
	hugo server --buildDrafts --verbose --source="${ROOT}/doc" --config="${ROOT}/doc/config.toml" --port=1377 --disableFastRender

docs_deploy:
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

clean:
	@echo MARK: clean
	rm -rf ${ROOT}/dist
	rm -rf ${ROOT}/build/public
	rm -rf ${ROOT}/static_source/admin/node_modules
	docker rmi -f $(docker images -aq)
	docker rm -vf $(docker ps -aq)

front_client:
	@echo MARK: generate front client lib
	npx swagger-typescript-api@12.0.4 --axios -p ./internal/api/api.swagger.yaml -o ./static_source/admin/src/api -n stub_new.ts

typedoc:
	@echo MARK: typedoc
	npx typedoc --tsconfig ./data/scripts/tsconfig.json --out ./internal/api/typedoc ./data/scripts/global.d.ts

.PHONY: build_darwin_arm64
build_darwin_arm64:
	@echo MARK: build local artefact
	RELEASE_VERSION=${RELEASE_VERSION} GO_BUILD_LDFLAGS=${GO_BUILD_LDFLAGS} docker buildx bake artifact-darwin-arm64

.PHONY: build_artifacts
build_artifacts:
	@echo MARK: build all artefacts
	RELEASE_VERSION=${RELEASE_VERSION} GO_BUILD_LDFLAGS=${GO_BUILD_LDFLAGS} docker buildx bake artifact-all

.PHONY: local_build
local_build:
	@echo MARK: local build
	RELEASE_VERSION=${RELEASE_VERSION} GO_BUILD_LDFLAGS=${GO_BUILD_LDFLAGS} docker buildx bake image-linux-arm64 --load

.PHONY: test_build
test_build:
	@echo MARK: local build
	RELEASE_VERSION=${RELEASE_VERSION} GO_BUILD_LDFLAGS=${GO_BUILD_LDFLAGS} docker buildx bake image-all

.PHONY: publish # Push the image to the remote registry
publish:
	@echo MARK: push docker image
	echo ${DOCKER_PASSWORD} | docker login -u ${DOCKER_USERNAME} --password-stdin && \
	RELEASE_VERSION=${RELEASE_VERSION} GO_BUILD_LDFLAGS=${GO_BUILD_LDFLAGS} docker buildx bake image-all --push

.PHONY: create_env
create_env:
	echo "RELEASE_VERSION=\"${RELEASE_VERSION}\"\nGO_BUILD_LDFLAGS=\"${GO_BUILD_LDFLAGS}\"" >> .env

.PHONY: build_archive
build_archive:
	cd ${ROOT}/dist/linux_amd64 && ls -l && tar -zcf ${ROOT}/linux_amd64.tar.gz .
	cd ${ROOT}/dist/linux_arm64 && ls -l && tar -zcf ${ROOT}/linux_arm64.tar.gz .
	cd ${ROOT}/dist/linux_arm_v6 && ls -l && tar -zcf ${ROOT}/linux_arm_v6.tar.gz .
	cd ${ROOT}/dist/linux_arm_v7 && ls -l && tar -zcf ${ROOT}/linux_arm_v7.tar.gz .
	cd ${ROOT}/dist/linux_ppc64le && ls -l && tar -zcf ${ROOT}/linux_ppc64le.tar.gz .
	cd ${ROOT}/dist/linux_riscv64 && ls -l && tar -zcf ${ROOT}/linux_riscv64.tar.gz .
	cd ${ROOT}/dist/linux_s390x && ls -l && tar -zcf ${ROOT}/linux_s390x.tar.gz .
	cd ${ROOT}/dist/windows_amd64 && ls -l && tar -zcf ${ROOT}/windows_amd64.tar.gz .
	cd ${ROOT}/dist/windows_arm64 && ls -l && tar -zcf ${ROOT}/windows_arm64.tar.gz .

.PHONY: build_artifact_public
build_artifact_public:
	@echo MARK: local build
	docker buildx bake artifact-public
