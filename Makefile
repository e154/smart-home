.PHONY: get_deps fmt
.DEFAULT_GOAL := build
tests: lint test
all: build_public build_linux_amd64 build_structure build_common_structure build_archive docker_image_linux_amd64
deploy: docker_image_upload

EXEC=server
CLI=cli
ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
SERVER_DIR = ${ROOT}/tmp/${EXEC}
COMMON_DIR = ${ROOT}/tmp/common
ARCHIVE=smart-home-common.tar.gz

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
RELEASE_VERSION ?= test
DOCKER_IMAGE_VER=${DOCKER_ACCOUNT}/${IMAGE}:${RELEASE_VERSION}
DOCKER_IMAGE_LATEST=${DOCKER_ACCOUNT}/${IMAGE}:latest

VERSION_VAR=${PROJECT}/version.VersionString
REV_VAR=${PROJECT}/version.RevisionString
REV_URL_VAR=${PROJECT}/version.RevisionURLString
GENERATED_VAR=${PROJECT}/version.GeneratedString
DEVELOPERS_VAR=${PROJECT}/version.DevelopersString
BUILD_NUMBER_VAR=${PROJECT}/version.BuildNumString
DOCKER_IMAGE_VAR=${PROJECT}/version.DockerImageString

GO_BUILD_LDFLAGS= -s -w -linkmode external -X ${VERSION_VAR}=${RELEASE_VERSION} -X ${REV_VAR}=${REV_VALUE} -X ${REV_URL_VAR}=${REV_URL_VALUE} -X ${GENERATED_VAR}=${GENERATED_VALUE} -X ${DEVELOPERS_VAR}=${DEVELOPERS_VALUE} -X ${BUILD_NUMBER_VAR}=${BUILD_NUMBER_VALUE} -X ${DOCKER_IMAGE_VAR}=${DOCKER_IMAGE_VER}
GO_BUILD_FLAGS= -a -installsuffix cgo -v --ldflags '${GO_BUILD_LDFLAGS}'
GO_BUILD_ENV=CGO_ENABLED=1
GO_BUILD_TAGS= -tags 'production'
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

# linux
build_linux_x86:
	@echo MARK: build linux x86
	rm -rf ${ROOT}/${EXEC}-linux-x86
	mkdir -p ${ROOT}/${EXEC}-linux-x86
	./bin/install_vosk.sh linux x86 ${ROOT}/${EXEC}-linux-x86
	export CGO_LDFLAGS_ALLOW="-Wl,-rpath.*" && \
	export CGO_CFLAGS="-I${ROOT}/${EXEC}-linux-x86" && \
	export CGO_LDFLAGS="-L${ROOT}/${EXEC}-linux-x86" && \
	export LD_LIBRARY_PATH="${ROOT}/${EXEC}-linux-x86:$LD_LIBRARY_PATH" && \
	export DYLD_LIBRARY_PATH="${ROOT}/${EXEC}-linux-x86" && \
	${GO_BUILD_ENV} CXX='g++ -m32' CC='gcc -m32' GOOS=linux GOARCH=386 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-linux-x86/server

build_linux_amd64:
	@echo MARK: build linux amd64
	rm -rf ${ROOT}/${EXEC}-linux-amd64
	mkdir -p ${ROOT}/${EXEC}-linux-amd64
	./bin/install_vosk.sh linux x86_64 ${ROOT}/${EXEC}-linux-amd64
	export CGO_LDFLAGS_ALLOW="-Wl,-rpath.*" && \
	export CGO_CFLAGS="-I${ROOT}/${EXEC}-linux-amd64" && \
	export CGO_LDFLAGS="-L${ROOT}/${EXEC}-linux-amd64" && \
	export LD_LIBRARY_PATH="${ROOT}/${EXEC}-linux-amd64:$LD_LIBRARY_PATH" && \
	export DYLD_LIBRARY_PATH="${ROOT}/${EXEC}-linux-amd64" && \
	${GO_BUILD_ENV} GOOS=linux GOARCH=amd64 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-linux-amd64/server

build_linux_armv6:
	@echo MARK: build linux armv6
	rm -rf ${ROOT}/${EXEC}-linux-arm-6
	mkdir -p ${ROOT}/${EXEC}-linux-arm-6
	./bin/install_vosk.sh linux armv7l ${ROOT}/${EXEC}-linux-arm-6
	export CGO_LDFLAGS_ALLOW="-Wl,-rpath.*" && \
	export CGO_CFLAGS="-I${ROOT}/${EXEC}-linux-arm-6" && \
	export CGO_LDFLAGS="-L${ROOT}/${EXEC}-linux-arm-6" && \
	export LD_LIBRARY_PATH="${ROOT}/${EXEC}-linux-arm-6:$LD_LIBRARY_PATH" && \
	export DYLD_LIBRARY_PATH="${ROOT}/${EXEC}-linux-arm-6" && \
	${GO_BUILD_ENV} CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-linux-arm-6/server

build_linux_armv7l:
	@echo MARK: build linux armv7l
	rm -rf ${ROOT}/${EXEC}-linux-arm-7
	mkdir -p ${ROOT}/${EXEC}-linux-arm-7
	./bin/install_vosk.sh linux armv7l ${ROOT}/${EXEC}-linux-arm-7
	export CGO_LDFLAGS_ALLOW="-Wl,-rpath.*" && \
	export CGO_CFLAGS="-I${ROOT}/${EXEC}-linux-arm-7" && \
	export CGO_LDFLAGS="-L${ROOT}/${EXEC}-linux-arm-7" && \
	export LD_LIBRARY_PATH="${ROOT}/${EXEC}-linux-arm-7:$LD_LIBRARY_PATH" && \
	export DYLD_LIBRARY_PATH="${ROOT}/${EXEC}-linux-arm-7" && \
	${GO_BUILD_ENV} CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=7 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-linux-arm-7/server

build_linux_arm64:
	@echo MARK: build linux arm64
	rm -rf ${ROOT}/${EXEC}-linux-arm64
	mkdir -p ${ROOT}/${EXEC}-linux-arm64
	./bin/install_vosk.sh linux aarch64 ${ROOT}/${EXEC}-linux-arm64
	export CGO_LDFLAGS_ALLOW="-Wl,-rpath.*" && \
	export CGO_CFLAGS="-I${ROOT}/${EXEC}-linux-arm64" && \
	export CGO_LDFLAGS="-L${ROOT}/${EXEC}-linux-arm64" && \
	export LD_LIBRARY_PATH="${ROOT}/${EXEC}-linux-arm64:$LD_LIBRARY_PATH" && \
	export DYLD_LIBRARY_PATH="${ROOT}/${EXEC}-linux-arm64" && \
	${GO_BUILD_ENV} CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-linux-arm64/server

# windows
build_windows_amd64:
	@echo MARK: build windows amd64
	rm -rf ${ROOT}/${EXEC}-windows-amd64
	mkdir -p ${ROOT}/${EXEC}-windows-amd64
	./bin/install_vosk.sh win win64 ${ROOT}/${EXEC}-windows-amd64
	export CGO_LDFLAGS_ALLOW="-Wl,-rpath.*" && \
	export CGO_CFLAGS="-I${ROOT}/${EXEC}-windows-amd64" && \
	export CGO_LDFLAGS="-L${ROOT}/${EXEC}-windows-amd64 -lvosk -L/usr/lib/x86_64-linux-gnu -ldl -lpthread" && \
	export LD_LIBRARY_PATH="${ROOT}/${EXEC}-windows-amd64:$LD_LIBRARY_PATH" && \
	export DYLD_LIBRARY_PATH="${ROOT}/${EXEC}-windows-amd64" && \
	${GO_BUILD_ENV} CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-windows-amd64/server.exe

build_windows_x86:
	@echo MARK: build windows x86
	rm -rf ${ROOT}/${EXEC}-windows-x86
	mkdir -p ${ROOT}/${EXEC}-windows-x86
	./bin/install_vosk.sh win win32 ${ROOT}/${EXEC}-windows-x86
	export CGO_LDFLAGS_ALLOW="-Wl,-rpath.*" && \
	export CGO_CFLAGS="-I${ROOT}/${EXEC}-windows-x86" && \
	export CGO_LDFLAGS="-L/usr/local/lib -L/usr/lib -L${ROOT}/${EXEC}-windows-x86 -lvosk -L/usr/lib/i386-linux-gnu -ldl -lpthread" && \
	export LD_LIBRARY_PATH="${ROOT}/${EXEC}-windows-x86:$LD_LIBRARY_PATH" && \
	export DYLD_LIBRARY_PATH="${ROOT}/${EXEC}-windows-x86" && \
	${GO_BUILD_ENV} CXX=i686-w64-mingw32-g++ CC=i686-w64-mingw32-gcc GOOS=windows GOARCH=386 go build -ldflags ${GO_BUILD_FLAGS} ${GO_BUILD_TAGS} -o ${ROOT}/${EXEC}-windows-x86/server.exe

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
	mkdir -p ${ROOT}/api/stub && \
	oapi-codegen -generate server -package stub ${ROOT}/api/api.swagger.yaml > ${ROOT}/api/stub/server.go && \
	oapi-codegen -generate types -package stub ${ROOT}/api/api.swagger.yaml > ${ROOT}/api/stub/types.go

build_structure:
	@echo MARK: create server structure
	mkdir -p ${SERVER_DIR}
	mkdir -p ${SERVER_DIR}/snapshots
	cd ${SERVER_DIR}
	cp -r ${ROOT}/conf ${SERVER_DIR}
	cp -r ${ROOT}/data ${SERVER_DIR}
	cp ${ROOT}/LICENSE ${SERVER_DIR}
	cp ${ROOT}/README* ${SERVER_DIR}
	cp ${ROOT}/CONTRIBUTING.md ${SERVER_DIR}
	chmod +x ${SERVER_DIR}/data/scripts/ping.sh

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

build_archive:
	@echo MARK: build app archive
	cd ${COMMON_DIR} && ls -l && tar -zcf ${ROOT}/${ARCHIVE} .

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

docker_image_linux_x86:
	echo ${ROOT}/${EXEC}-linux-x86
	cd ${ROOT}/${EXEC}-linux-x86 && ls -ll && docker build --platform linux/386 -f ${ROOT}/bin/docker/Dockerfile -t ${DOCKER_ACCOUNT}/${IMAGE}-x86:${RELEASE_VERSION} .

docker_image_linux_amd64:
	echo ${ROOT}/${EXEC}-linux-amd64
	cd ${ROOT}/${EXEC}-linux-amd64 && ls -ll && docker build --platform linux/amd64 -f ${ROOT}/bin/docker/Dockerfile -t ${DOCKER_ACCOUNT}/${IMAGE}-amd64:${RELEASE_VERSION} .

docker_image_linux_armv7l:
	echo ${ROOT}/${EXEC}-linux-arm-7
	cd ${ROOT}/${EXEC}-linux-arm-7 && ls -ll && docker build --platform linux/arm/v7 -f ${ROOT}/bin/docker/Dockerfile -t ${DOCKER_ACCOUNT}/${IMAGE}-arm-7:${RELEASE_VERSION} .

docker_image_linux_arm64:
	echo ${ROOT}/${EXEC}-linux-arm64
	cd ${ROOT}/${EXEC}-linux-arm64 && ls -ll && docker build --platform linux/arm64 -f ${ROOT}/bin/docker/Dockerfile -t ${DOCKER_ACCOUNT}/${IMAGE}-arm64:${RELEASE_VERSION} .

docker_image_upload:
	echo ${DOCKER_PASSWORD} | docker login -u ${DOCKER_USERNAME} --password-stdin
	docker push ${DOCKER_ACCOUNT}/${IMAGE}-amd64:${RELEASE_VERSION}
	docker push ${DOCKER_ACCOUNT}/${IMAGE}-arm64:${RELEASE_VERSION}
	docker manifest create --amend ${DOCKER_ACCOUNT}/${IMAGE}:${RELEASE_VERSION} \
	${DOCKER_ACCOUNT}/${IMAGE}-amd64:${RELEASE_VERSION} \
	${DOCKER_ACCOUNT}/${IMAGE}-arm64:${RELEASE_VERSION}
	docker manifest push {DOCKER_ACCOUNT}/${IMAGE}:${RELEASE_VERSION}

clean:
	@echo MARK: clean
	rm -rf ${SERVER_DIR}
	rm -f ${ROOT}/${EXEC}-*
	rm -f ${ROOT}/${CLI}-*
	rm -f ${ROOT}/libvosk*
	rm -f ${HOME}/${ARCHIVE}

front_client:
	@echo MARK: generate front client lib
	npx swagger-typescript-api@12.0.4 --axios -p ./api/api.swagger.yaml -o ./static_source/admin/src/api -n stub_new.ts

typedoc:
	@echo MARK: typedoc
	npx typedoc --tsconfig ./data/scripts/tsconfig.json --out ./api/typedoc ./data/scripts/global.d.ts
