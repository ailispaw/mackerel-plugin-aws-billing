FROM alpine:3.5

ENV GOPATH="/tmp/go"

RUN apk --no-cache add --virtual build-deps go git build-base

COPY . "$GOPATH/src/github.com/littlekbt/mackerel-plugin-aws-billing"

RUN cd "$GOPATH/src/github.com/littlekbt/mackerel-plugin-aws-billing" && \
    go get -d -v ./... && \
    go build -o /usr/bin/mackerel-plugin-aws-billing -ldflags "-w -s"

ENTRYPOINT [ "mackerel-plugin-aws-billing" ]
