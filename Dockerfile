FROM golang:1.22.2-alpine as builder


RUN apk update && \
    apk add --no-cache git &&\
    apk add --no-cache curl

FROM builder as local

WORKDIR /app

RUN go install github.com/go-task/task/v3/cmd/task@latest && \
    go install github.com/githubnemo/CompileDaemon@latest

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download


FROM local as dev_service

ARG service_port=8080
ARG metric_port_service=8081

EXPOSE ${service_port}
EXPOSE ${metric_port_service}



FROM local as dev_consumer

ARG metric_port_consumer=8082

EXPOSE ${metric_port_consumer}



FROM builder as service_builder

MAINTAINER go-feast

COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest

RUN ./task build-api-server

FROM alpine:latest as prod_service

WORKDIR /app

COPY --from=service_builder /app/bin/api-server .

ARG service_port
ARG service_metrics_port

EXPOSE ${service_port}
EXPOSE ${service_metrics_port}

CMD ["./api-server"]

FROM builder as consumer_builder

MAINTAINER go-feast

COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest

RUN ./task build-api-consumer

FROM alpine:latest as prod_consumer

WORKDIR /app

COPY --from=consumer_builder /app/bin/api-consumer .

ARG consumer_metrics_port

EXPOSE ${consumer_metrics_port}

CMD ["./api-consumer"]



