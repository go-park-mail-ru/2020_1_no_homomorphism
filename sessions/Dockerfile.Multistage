FROM golang:1.14 AS stage1
ENV GO111MODULE=on
WORKDIR /go/src/stage1
COPY sessions /go/src/stage1
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine
WORKDIR /app
COPY --from=stage1  /go/src/stage1 /app
RUN chmod +x .
EXPOSE 8083/tcp
ENTRYPOINT [ "/app/no_homo_sessions" ]