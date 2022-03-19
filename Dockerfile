FROM golang:1.15.12 AS builder

ADD . /data
WORKDIR /data
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o srv


FROM alpine:3.6 AS runtime

RUN mkdir /usr/srv
COPY --from=builder /data/srv /usr/srv/srv
WORKDIR /usr/srv

ENV salt jqh8i912980j1rf1908wdj183
ENV LANG C.UTF-8
ENTRYPOINT ["./srv"]
