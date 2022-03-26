FROM golang:1.18 as build
RUN update-ca-certificates
MAINTAINER Filippov Alex <support@e154.ru>

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
ENTRYPOINT ["/server"]
USER nobody
ADD . /

EXPOSE 3000
EXPOSE 3001
EXPOSE 3002
EXPOSE 3003
EXPOSE 1883

ENV SERVER_HOST="0.0.0.0"
ENV SERVER_PORT="3000"
ENV PG_USER="smart_home"
ENV PG_PASS="smart_home"
ENV PG_HOST="postgres"
ENV PG_NAME="smart_home"
ENV PG_PORT="5432"
ENV PG_DEBUG="false"
ENV MODE="release"
ENV AUTO_MIGRATE="true"

VOLUME $APP_DIR/snapshots
VOLUME $APP_DIR/data
VOLUME $APP_DIR/conf
