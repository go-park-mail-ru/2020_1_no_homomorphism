FROM golang:1.14

COPY ./sessions go/src/session
RUN go install session

EXPOSE 8081/tcp

CMD ["session"]