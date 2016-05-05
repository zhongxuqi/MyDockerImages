FROM golang:latest

COPY src $GOPATH/src
COPY bin $GOPATH/bin

EXPOSE 8089 37077 8090 37078

ENV PATH .:$PATH

WORKDIR $GOPATH/bin

RUN go build DesertEagleSite
RUN go build ServiceSite

CMD ["./run.sh"]

