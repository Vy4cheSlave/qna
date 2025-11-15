FROM golang:1.25.3-alpine AS builder

ENV CGO_ENABLED 0

ENV GOOS linux

ENV GOPROXY direct

RUN apk update --no-cache && apk add --no-cache ca-certificates git && apk add --no-cache --upgrade bash

WORKDIR /builder

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build cmd/main.go

FROM alpine:latest AS runner

RUN apk update --no-cache && apk add --no-cache ca-certificates && apk add --no-cache --upgrade bash

WORKDIR /build

COPY . .

COPY --from=builder /builder/main /build

CMD ["./main"]