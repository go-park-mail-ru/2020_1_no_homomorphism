FROM golang:1.14 AS no_homo_main_1
ENV GO111MODULE=on
WORKDIR /go/src/no_homo_main_1
COPY main /go/src/no_homo_main_1
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 FILE_ROOT="resources" FILE_SERVER="https://localhost:8083" go build .

FROM alpine
WORKDIR /app
COPY --from=no_homo_main_1  /go/src/no_homo_main_1 /app
RUN chmod +x  .
EXPOSE 8081/tcp
EXPOSE 5432/tcp
ENTRYPOINT [ "/app/no_homo_main" ]