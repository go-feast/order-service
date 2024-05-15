# Build stage
FROM golang:1.22.2-alpine AS builder

RUN apk update && \
    apk add --no-cache git &&\
    apk add --no-cache curl


FROM builder AS dev
# Set the working directory inside the container
WORKDIR /app
#metrics port
EXPOSE 8081

RUN go install github.com/go-task/task/v3/cmd/task@latest && \
    go install github.com/githubnemo/CompileDaemon@latest

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download