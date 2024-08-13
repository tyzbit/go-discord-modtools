FROM golang:1.23-alpine as build
WORKDIR /
COPY . ./

RUN apk add \
    build-base \
    git \
&&  go build -ldflags="-s -w"

FROM alpine
ENV GIN_MODE=release

COPY --from=build /go-discord-modtools /

CMD ["/go-discord-modtools"]
