FROM jerson/go:1.14 AS builder

ENV WORKDIR /app
ENV GOPROXY https://goproxy.io
WORKDIR ${WORKDIR}

COPY go.mod go.sum Makefile ./
RUN make deps

USER root

COPY config.toml-dist config.toml
COPY . .

RUN make build

FROM jerson/base:1.3

LABEL maintainer="jeral17@gmail.com"

ENV BUILDER_PATH /app
ENV WORKDIR /app
WORKDIR ${WORKDIR}

COPY --from=builder ${BUILDER_PATH}/config.toml .
COPY --from=builder ${BUILDER_PATH}/nats-webhook .

EXPOSE 8080

ENTRYPOINT ["/app/nats-webhook"]