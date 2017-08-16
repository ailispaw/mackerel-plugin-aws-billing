FROM alpine:3.5

ARG VERSION
ENV GOPATH="/tmp/go"

RUN apk --no-cache --update add ca-certificates \
 && rm -rf /tmp/* /var/cache/apk/*

RUN apk --no-cache add --virtual build-deps git go build-base \
    \
 && mkdir -p "$GOPATH/src/github.com/littlekbt" \
 && cd "$GOPATH/src/github.com/littlekbt" \
 && git clone --depth 1 --branch $VERSION https://github.com/ailispaw/mackerel-plugin-aws-billing \
    \
 && cd mackerel-plugin-aws-billing/src \
 && go get -d -v ./... \
 && go build -o /usr/bin/mackerel-plugin-aws-billing -ldflags "-w -s" \
    \
 && apk del build-deps \
 && rm -rf /usr/lib/go \
 && rm -rf "${GOPATH}" \
 && rm -rf /tmp/* /var/cache/apk/*

ENTRYPOINT [ "mackerel-plugin-aws-billing" ]
