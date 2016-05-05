FROM golang:latest

COPY src $GOPATH/src
COPY bin $GOPATH/bin

EXPOSE 8089 37077

ENV PATH .:$PATH

WORKDIR $GOPATH/bin

RUN go build DesertEagleSite

CMD ["./run.sh"]

