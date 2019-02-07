FROM golang:alpine AS builder

RUN apk add git && \
    mkdir /root/bin

WORKDIR $GOPATH/src/proxy/
ADD src/proxy/proxy.go .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /root/bin/proxy

WORKDIR $GOPATH/src/iptvgenerator/
ADD src/iptvgenerator/iptvgenerator.go .
RUN go get github.com/PuerkitoBio/goquery && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /root/bin/iptvgenerator

FROM alpine:latest

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=builder /root/bin/ /usr/bin

RUN echo "0 */3 * * * iptvgenerator > /dev/stdout" >> /etc/crontabs/root
VOLUME [ "/data" ]

EXPOSE 8080

CMD ["sh", "-c", "crond && proxy"]