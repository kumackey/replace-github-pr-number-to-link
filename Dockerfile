# Use the official Golang image as the base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy the Go source file(s) into the container
COPY main.go .

# Build the Go application
RUN go build -o /usr/local/bin/app main.go

# Configure the container to run the Go application
ENTRYPOINT ["/usr/local/bin/app"]
