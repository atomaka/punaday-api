FROM alpine:latest
COPY . /go/src/github.com/atomaka/punaday-api
RUN apk update \
    && apk add --no-cache go git \
    && cd /go/src/github.com/atomaka/punaday-api \
    && export GOPATH=/go \
    && go get \
    && go build -o /bin/punaday-api \
    && rm -rf /go \
    && apk del --purge git go \
    && rm -rf /var/cache/apk*
ENTRYPOINT ["/bin/punaday-api"]
