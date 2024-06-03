FROM golang:1.22.2-alpine as builder

MAINTAINER go-feast

WORKDIR /app

ARG service_port
ARG metric_port_service
ARG metric_port_consumer

RUN --mount=type=cache,target=/var/cache/apt\
    apk update && \
    apk add --no-cache git &&\
    apk add --no-cache curl

RUN --mount=type=cache,target=/var/cache/go/bin\
    go install github.com/go-task/task/v3/cmd/task@latest

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download the dependencies
RUN  --mount=type=cache,target=/var/cache/go/pkg/mod\
     go mod download

FROM builder as local

RUN --mount=type=cache,target=/var/cache/go/bin\
    go install github.com/githubnemo/CompileDaemon@latest

FROM local as dev_service

ARG service_port
ARG metric_port_service

EXPOSE ${service_port}
EXPOSE ${metric_port_service}

FROM local as dev_consumer

ARG metric_port_consumer

EXPOSE ${metric_port_consumer}

FROM builder as service_builder

# Copy the rest of the application source code
COPY . .

RUN /go/bin/task build-api-server

FROM alpine:latest as prod_service

WORKDIR /app

ARG service_port
ARG service_metrics_port

COPY --from=service_builder /app/bin/api-server .


EXPOSE ${service_port}
EXPOSE ${service_metrics_port}

CMD ["./api-server"]

FROM builder as consumer_builder

# Copy the rest of the application source code
COPY . .

RUN /go/bin/task build-api-consumer

FROM alpine:latest as prod_consumer

WORKDIR /app

ARG consumer_metrics_port

COPY --from=consumer_builder /app/bin/api-consumer .

EXPOSE ${consumer_metrics_port}

CMD ["./api-consumer"]