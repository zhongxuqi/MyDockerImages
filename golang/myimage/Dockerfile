FROM golang:latest

COPY src /go/src

EXPOSE 8089 37077
# EXPOSE 8090
# EXPOSE 7077

ENV PATH .:$PATH

WORKDIR $GOPATH/src/DesertEagleSite

RUN go build DesertEagleSite

CMD ["./DesertEagleSite"]
