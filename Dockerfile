FROM golang

ADD ./service /go/src/github.com/macchid/questandans

RUN go install github.com/macchid/questandans

ENTRYPOINT [ /go/bin/questandans ]

EXPOSE 8080

