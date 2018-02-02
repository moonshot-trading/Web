FROM golang:1.9-alpine3.6

RUN mkdir -p /go/src/github.com/moonshot-trading/Web/

ADD . /go/src/github.com/moonshot-trading/Web/

RUN go get github.com/moonshot-trading/Web/server
RUN go install github.com/moonshot-trading/Web/server

ENTRYPOINT /go/bin/server

EXPOSE 8080
