FROM golang:latest

COPY src /go/src

EXPOSE 8089 37077
# EXPOSE 8090
# EXPOSE 7077

ENV PATH ".:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

WORKDIR $GOPATH/src/DesertEagleSite

CMD ["./DesertEagleSite"]
