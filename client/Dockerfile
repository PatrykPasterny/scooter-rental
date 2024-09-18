# Use a base image with Go installed
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the entire project to the container
COPY . .

# Build the Go binary
RUN go build -o myapp

# Specify the command to run your app when the container starts
CMD ["./myapp"]