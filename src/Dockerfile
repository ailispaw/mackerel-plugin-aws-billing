FROM alpine:3.5

ENV GOPATH="/tmp/go"

COPY . "$GOPATH/src/github.com/littlekbt/mackerel-plugin-aws-billing/src"

RUN apk --no-cache add --virtual build-deps go git build-base && \
    cd "$GOPATH/src/github.com/littlekbt/mackerel-plugin-aws-billing/src" && \
    go get -d -v ./... && \
    go build -o /usr/bin/mackerel-plugin-aws-billing -ldflags "-w -s" && \
    apk del build-deps && \
    rm -rf /usr/lib/go && \
    rm -rf "${GOPATH}"

RUN apk --no-cache add ca-certificates

ENTRYPOINT [ "mackerel-plugin-aws-billing" ]
