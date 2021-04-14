FROM golang:1.16.3-alpine3.12 as builder
RUN mkdir /build
ADD . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o janitor cmd/janitor/*.go

FROM alpine:3.10.8
RUN apk add --no-cache bash
COPY .github/actions/janitor/entrypoint.sh /entrypoint.sh
COPY --from=builder /build/janitor /janitor
RUN chmod +x /janitor

ENTRYPOINT ["/entrypoint.sh"]