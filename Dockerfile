# Use official Go image as builder
FROM golang:1.24

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the app
RUN go build -o weather-api

# Run the app
CMD ["./weather-api"]