FROM golang:1.20-alpine AS build

WORKDIR /workspace

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o store-server ./cmd/store/main.go

FROM alpine:latest

WORKDIR /store

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    mkdir -p /opt/store/bin && \
    mkdir -p /opt/store/conf && \
    mkdir -p /opt/store/log

COPY --from=build /workspace/store-server /opt/store/bin/

VOLUME /opt/store/conf /opt/store/log

ENTRYPOINT ["/opt/store/bin/store-server"]
CMD ["-c", "/opt/store/conf/store-server.yaml"]
