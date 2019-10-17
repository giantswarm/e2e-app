FROM alpine:3.10

RUN apk add --update ca-certificates \
    && rm -rf /var/cache/apk/*

ADD ./e2e-app /e2e-app

ENTRYPOINT ["/e2e-app"]
