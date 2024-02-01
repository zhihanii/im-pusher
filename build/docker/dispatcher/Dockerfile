# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS build

WORKDIR /workspace

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o dispatcher ./cmd/dispatcher/main.go

FROM alpine:latest

WORKDIR /dispatcher

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    mkdir -p /opt/dispatcher/bin && \
    mkdir -p /opt/dispatcher/conf && \
    mkdir -p /opt/dispatcher/log

COPY --from=build /workspace/dispatcher /opt/dispatcher/bin/

VOLUME /opt/dispatcher/conf /opt/dispatcher/log

EXPOSE 8103 8203

ENTRYPOINT ["/opt/dispatcher/bin/dispatcher"]
CMD ["-c", "/opt/dispatcher/conf/dispatcher.yaml"]
