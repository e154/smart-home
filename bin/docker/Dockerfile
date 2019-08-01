FROM debian
MAINTAINER Filippov Alex <af@suretly.com>

ENV APP_DIR /opt/smart-home-server
RUN sed 's/main$/main universe/' -i /etc/apt/sources.list \
    && apt-get update \
    && apt-get upgrade -y \
    && apt-get install ca-certificates -y \
    && mkdir -p $APP_DIR

WORKDIR $APP_DIR

ADD . $APP_DIR

ENTRYPOINT (cd $APP_DIR && ./server)

EXPOSE 3000
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