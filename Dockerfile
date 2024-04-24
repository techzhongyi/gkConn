FROM golang:1.17-buster as build
WORKDIR /app
ENV GOPROXY="https://goproxy.cn,direct"
COPY comlibgo ./comlibgo
WORKDIR /app/gkConn
COPY gkConn/go.mod .
COPY gkConn/go.sum .

RUN go mod download

COPY gkConn/src ./src
COPY gkConn/config.yaml .
COPY gkConn/startup.go .

RUN CGO_ENABLED=0 go build -o docker-gkConn

## Deploy
FROM alpine:3.9
RUN apk add ca-certificates
WORKDIR /
COPY --from=build /app/gkConn/docker-gkConn /app/exec/docker-gkConn
COPY --from=build /app/gkConn/config.yaml /app/exec/config.yaml
COPY --from=build /app/gkConn/__all_apis/ /app/

WORKDIR /app/exec

ENTRYPOINT ["/app/exec/docker-gkConn"]

