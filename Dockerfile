FROM golang:1.24.4-bookworm@sha256:97162678719a516c12d5fb4b08266ab04802358cff63697ab1584be29ee8995c AS builder-backend

COPY go.* .
COPY *.go .

RUN go build .

RUN ls -al

FROM debian:bookworm-slim

RUN mkdir -p /app

WORKDIR /app

COPY LICENSE .
COPY --from=builder-backend /go/inittmpl .
ENV PATH="$PATH:/app"

CMD ["inittmpl"]