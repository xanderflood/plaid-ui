FROM alpine:3.10

RUN apk add --no-cache ca-certificates

COPY api api
COPY templates templates
COPY static static

CMD ["./api"]
