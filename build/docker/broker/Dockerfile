FROM golang:1.20-alpine AS build

WORKDIR /workspace

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o broker ./cmd/broker/main.go

FROM alpine:latest

WORKDIR /broker

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    mkdir -p /opt/broker/bin && \
    mkdir -p /opt/broker/conf && \
    mkdir -p /opt/broker/log && \
    mkdir -p /opt/broker/pprof

COPY --from=build /workspace/broker /opt/broker/bin/

VOLUME /opt/broker/conf /opt/broker/log /opt/broker/pprof

EXPOSE 8101 8201 8082

ENTRYPOINT ["/opt/broker/bin/broker"]
CMD ["-c", "/opt/broker/conf/broker.yaml"]
