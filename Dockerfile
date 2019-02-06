FROM golang:alpine AS builder

RUN apk update && \
    apk upgrade && \
    apk add git && \
    mkdir /root/bin

WORKDIR $GOPATH/src/proxy/
ADD src/proxy/proxy.go .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /root/bin/proxy

WORKDIR $GOPATH/src/iptvgenerator/
ADD src/iptvgenerator/iptvgenerator.go .
RUN go get github.com/tebeka/selenium && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /root/bin/iptvgenerator


FROM alpine:latest

COPY --from=builder /root/bin/ /usr/bin

RUN apk update && \
    apk upgrade && \
    apk add chromium chromium-chromedriver xvfb --no-cache && \
    echo "0 5 * * * iptvgenerator > /dev/stdout" >> /etc/crontabs/root

VOLUME [ "/data" ]

EXPOSE 8080

CMD ["sh", "-c", "crond && proxy"]