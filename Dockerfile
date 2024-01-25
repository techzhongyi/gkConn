FROM golang:1.17-buster as build
WORKDIR /app
ENV GOPROXY="https://goproxy.cn,direct"
COPY comlibgo ./comlibgo
WORKDIR /app/httpCore
COPY httpCore/go.mod .
COPY httpCore/go.sum .

RUN go mod download

COPY httpCore/src ./src
COPY httpCore/config.yaml .
COPY httpCore/startup.go .

COPY httpCore/__all_apis ./__all_apis

RUN CGO_ENABLED=0 go build -o docker-httpCore

## Deploy
FROM alpine:3.9
RUN apk add ca-certificates
WORKDIR /
COPY --from=build /app/httpCore/docker-httpCore /app/exec/docker-httpCore
COPY --from=build /app/httpCore/config.yaml /app/exec/config.yaml
COPY --from=build /app/httpCore/__all_apis/ /app/

WORKDIR /app/exec

ENTRYPOINT ["/app/exec/docker-httpCore"]

