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

ADD src/bin /usr/bin
COPY --from=builder /root/bin/ /usr/bin

RUN apk update && \
    apk upgrade && \
    apk add chromium chromium-chromedriver xvfb openjdk8-jre --no-cache && \
    echo "@reboot proxy" >> /etc/crontabs/root && \
    echo "@reboot iptvgenerator" >> /etc/crontabs/root && \
    echo "50 4 * * * iptvgenerator" >> /etc/crontabs/root

EXPOSE 8080

CMD ["crond", "-f"]