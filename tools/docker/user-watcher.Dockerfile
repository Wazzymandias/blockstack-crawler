FROM golang:1.10.3-stretch as builder

ENV DST_DIR $GOPATH/src/github.com/Wazzymandias/blockstack-crawler

ENV GLIDE_URL github.com/Masterminds/glide

RUN mkdir -p ${DST_DIR}
RUN go get -u $GLIDE_URL

COPY . ${DST_DIR}

RUN cd ${DST_DIR} && \
    glide install && \
    go install -v

RUN unset DST_DIR

FROM debian:stable-slim

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get -y install cron

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

ENV CRAWLER_BIN_DIR /usr/local/bin

COPY --from=builder $GOPATH/bin/blockstack-crawler ${CRAWLER_BIN_DIR}

ADD tools/cron/watcher-crontab /etc/cron.d/watcher-cron
RUN chmod 0644 /etc/cron.d/watcher-cron
RUN touch /var/log/cron.log

COPY tools/dumb-init /usr/bin/dumb-init

RUN cat /etc/pam.d/cron |sed -e "s/required     pam_loginuid.so/optional     pam_loginuid.so/g" > /tmp/cron && mv /tmp/cron /etc/pam.d/cron

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["cron", "-f"]
