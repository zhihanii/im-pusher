FROM golang:1.20-alpine AS build

WORKDIR /workspace

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o logic ./cmd/logic/main.go

FROM alpine:latest

WORKDIR /logic

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    mkdir -p /opt/logic/bin && \
    mkdir -p /opt/logic/conf && \
    mkdir -p /opt/logic/log && \
    mkdir -p /opt/logic/pprof

COPY --from=build /workspace/logic /opt/logic/bin/

VOLUME /opt/logic/conf /opt/logic/log /opt/logic/pprof

EXPOSE 8102 8202

ENTRYPOINT ["/opt/logic/bin/logic"]
CMD ["-c", "/opt/logic/conf/logic.yaml"]
