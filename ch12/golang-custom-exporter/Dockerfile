# Use the official Golang base image
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy information about the used go modules
COPY go.mod .
COPY go.sum .

# Only download updates if modules files have changed
RUN go mod download

# Copy the application source code
COPY . .

# Build the Go application
RUN go build -o main -ldflags '-linkmode external -w -extldflags "-static"'

# Bundle Stage
FROM alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the executable from the builder stage
COPY --from=builder /app/main .

# Set the command to run the executable
CMD ["/app/main"]
