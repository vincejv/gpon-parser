FROM library/golang:1.21 as build-env

RUN apt-get install -yq --no-install-recommends git

# Copy source + vendor
COPY . /go/src/github.com/vincejv/gpon-parser
WORKDIR /go/src/github.com/vincejv/gpon-parser

# Build
ENV GOPATH=/go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -v -a -ldflags "-s -w" -o /go/bin/gpon-parser .

FROM alpine:latest
RUN apk --update --no-cache add curl && rm -rf /var/cache/apk/*
COPY --from=build-env /go/bin/gpon-parser /usr/bin/gpon-parser
ENTRYPOINT ["gpon-parser"]

EXPOSE 8092