# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM crazymax/goxx:latest AS goxx
FROM --platform=$BUILDPLATFORM crazymax/osxcross:11.3 AS osxcross

FROM goxx AS base
ENV GO111MODULE=auto
ENV CGO_ENABLED=1
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,target=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS build
ARG TARGETPLATFORM
ARG RELEASE_VERSION
ARG GO_BUILD_LDFLAGS
RUN --mount=type=cache,sharing=private,target=/var/cache/apt \
  --mount=type=cache,sharing=private,target=/var/lib/apt/lists \
  goxx-apt-get install -y binutils gcc g++ pkg-config
RUN --mount=type=bind,source=.,rw \
  --mount=from=osxcross,target=/osxcross,src=/osxcross,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod <<EOT
BUILDMODE=
if [ "$(. goxx-env && echo $GOOS)" != "windows" ]; then
  case "$(. goxx-env && echo $GOARCH)" in
    mips*|ppc64)
      # pie build mode is not supported on mips architectures
      ;;
    *)
      BUILDMODE="-buildmode=pie"
      ;;
  esac
fi
LDFLAGS="$GO_BUILD_LDFLAGS -s -w"
if [ "$(. goxx-env && echo $GOOS)" = "linux" ]; then
  LDFLAGS="$LDFLAGS -extldflags -static"
fi
GO_BUILD_TAGS="-tags production,netgo,osusergo"
goxx-go env
goxx-go build -a -v -o /out/server -trimpath -ldflags "$LDFLAGS" $BUILDMODE $GO_BUILD_TAGS .
#mkdir -p /out/deps
#ls -R /out
#ldd-copy-dependencies.sh -b /out/server -t /out/deps
#find /usr -name "libdl*" | xargs -I % sh -c 'mkdir -p $(dirname /out/deps%); cp % /out/deps%;'
EOT

FROM scratch AS artifact
COPY --from=build /out /

FROM --platform=$BUILDPLATFORM debian:bookworm-slim
RUN apt-get update; \
    apt-get install -y --no-install-recommends \
      libpq5 \
      ca-certificates \
      postgresql-client-15 \
      iputils-ping; \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN update-ca-certificates

RUN mkdir -p /app
RUN chown nobody /app

WORKDIR /app/

COPY conf conf/
COPY data data/
COPY LICENSE .
COPY README.md .
COPY CONTRIBUTING.md .

COPY --from=build /out/server .

ENTRYPOINT [ "/app/server" ]

EXPOSE 3000
EXPOSE 3001
EXPOSE 3002
EXPOSE 3003
EXPOSE 1883
EXPOSE 8080
EXPOSE 8443

ENV PG_USER="smart_home"
ENV PG_PASS="smart_home"
ENV PG_HOST="postgres"
ENV PG_NAME="smart_home"
ENV PG_PORT="5432"
ENV PG_DEBUG="false"
ENV PG_MAX_IDLE_CONNS="10"
ENV PG_MAX_OPEN_CONNS="50"
ENV PG_CONN_MAX_LIFE_TIME="30"
ENV AUTO_MIGRATE="true"
ENV SNAPSHOT_DIR=""
ENV MODE="release"
ENV MQTT_PORT="1883"
ENV MQTT_RETRY_INTERVAL="20"
ENV MQTT_RETRY_CHECK_INTERVAL="20"
ENV MQTT_SESSION_EXPIRY_INTERVAL="0"
ENV MQTT_SESSION_EXPIRE_CHECK_INTERVAL="0"
ENV MQTT_QUEUE_QOS_0_MESSAGES="true"
ENV MQTT_MAX_INFLIGHT="32"
ENV MQTT_MAX_AWAIT_REL="100"
ENV MQTT_MAX_MSG_QUEUE="1000"
ENV MQTT_DELIVER_MODE="1"
ENV LOGGING="true"
ENV COLORED_LOGGING="false"
ENV ALEXA_HOST=""
ENV ALEXA_PORT="3003"
ENV API_HTTP_PORT="3001"
ENV API_HTTPS_PORT="3002"
ENV API_SWAGGER="true"
ENV API_GZIP="false"
ENV LANG="EN"
ENV ROOT_MODE="true"
ENV DOMAIN="localhost"
ENV PPROF="false"
ENV ROOT_SECRET=""
ENV PATH="$PATH:/app"

ENV GATE_API_HTTP_PORT="8080"
ENV GATE_API_HTTPS_PORT="8443"
ENV GATE_API_DEBUG="false"
ENV GATE_API_GZIP="true"
ENV GATE_DOMAIN="localhost"
ENV GATE_PPROF="false"
ENV GATE_HTTPS="false"
ENV GATE_PROXY_TIMEOUT="5"
ENV GATE_PROXY_IDLE_TIMEOUT="10"
ENV GATE_PROXY_SECRET_KEY=""

VOLUME /app/snapshots
VOLUME /app/data
VOLUME /app/conf
