FROM golang:1.15.12 AS builder

ADD . /data
WORKDIR /data
# COPY .ssh /root/.ssh

# RUN apt-get install -y git
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o srv


FROM alpine:3.6 AS runtime

RUN mkdir /usr/srv
COPY --from=builder /data/srv /usr/srv/srv
WORKDIR /usr/srv

ENV LANG C.UTF-8
ENTRYPOINT ["./srv"]
