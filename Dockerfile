FROM golang:1.22.2-alpine as builder

MAINTAINER go-feast

WORKDIR /app

ARG service_port
ARG metrics_port

RUN --mount=type=cache,target=/var/cache/apt\
    apk update && \
    apk add --no-cache git &&\
    apk add --no-cache curl

RUN --mount=type=cache,target=/var/cache/go/bin\
    go install github.com/go-task/task/v3/cmd/task@latest

FROM builder as local

RUN --mount=type=cache,target=/var/cache/go/bin\
    go install github.com/githubnemo/CompileDaemon@latest \

FROM builder as production_build

ARG build

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download the dependencies
RUN  --mount=type=cache,target=/var/cache/go/pkg/mod\
     go mod download

# Copy the rest of the application source code
COPY . .

RUN /go/bin/task ${build}

FROM alpine:latest as production

WORKDIR /app

ARG build
ARG service_port
ARG metrics_port

COPY --from=production_build /app/bin/${build} .


EXPOSE ${service_port}
EXPOSE ${metrics_port}

CMD ["./${build}"]
