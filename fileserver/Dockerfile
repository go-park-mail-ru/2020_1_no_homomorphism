FROM golang:1.14 AS stage1
WORKDIR /go/src/stage1
COPY fileserver /go/src/stage1
COPY resources /go/src/stage1/resources
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine
WORKDIR /app
COPY --from=stage1  /go/src/stage1 /app
RUN chmod +x ./fileserver
EXPOSE 8082
ENTRYPOINT [ "/app/fileserver" ]
