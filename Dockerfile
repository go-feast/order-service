# Build Stage
FROM golang:1.22.2 AS build
MAINTAINER go-feast

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files for efficient caching
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

RUN ./install-task.sh

RUN ./task build

# Production Stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the build stage
COPY --from=build /app/bin/api-server .

# Expose the port that the application listens on
EXPOSE 8080

# Maybe expose port foe metrics
# EXPOSE 8081

# Command to run the application
CMD ["./api-server"]
