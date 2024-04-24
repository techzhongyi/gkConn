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

## Deploy..
FROM alpine:3.9
RUN apk add ca-certificates
WORKDIR /

COPY --from=build /app/gkConn/docker-gkConn /app/docker-gkConn
COPY --from=build /app/gkConn/config.yaml /app/config.yaml

RUN ls /app

ENTRYPOINT ["/app/docker-gkConn"]


# 注意！！！ 该docker file 必须在上级目录运行 /app/comligo  /app/patchCore  /app/Dockerfile  这种目录结构
# 通过Jenkins部署实现上述自动化打包构建
