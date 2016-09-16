FROM golang:1.7.1-alpine

RUN apk update && apk upgrade && apk add git

RUN mkdir -p /go/src/github.com/byuoitav
ADD . /go/src/github.com/byuoitav/london-audio-microservice

WORKDIR /go/src/github.com/byuoitav/london-audio-microservice
RUN go get -d -v
RUN go install -v

CMD ["/go/bin/london-audio-microservice"]

EXPOSE 8009
