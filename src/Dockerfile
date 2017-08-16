FROM alpine:3.5

ENV GOPATH="/tmp/go"

RUN apk --no-cache --update add ca-certificates \
 && rm -rf /tmp/* /var/cache/apk/*

RUN apk --no-cache --update add --virtual build-deps go git build-base

COPY . "$GOPATH/src/github.com/littlekbt/mackerel-plugin-aws-billing/src"

RUN cd "$GOPATH/src/github.com/littlekbt/mackerel-plugin-aws-billing/src" \
 && go get -d -v ./... \
 && go build -o /usr/bin/mackerel-plugin-aws-billing -ldflags "-w -s" \
    \
 && apk del build-deps \
 && rm -rf /usr/lib/go \
 && rm -rf "${GOPATH}" \
 && rm -rf /tmp/* /var/cache/apk/*

ENTRYPOINT [ "mackerel-plugin-aws-billing" ]
