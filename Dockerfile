# Step 1: Build the Go binary
FROM golang:1.21-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install CA certificates in the builder stage (useful for any external package downloads)
RUN apk add --no-cache ca-certificates

WORKDIR /build

# Copy the necessary files to the build directory
COPY go.mod go.sum *.go ./
COPY common ./common
COPY model ./model
COPY s3 ./s3
COPY *.json ./
COPY *.properties ./

# List files to verify correct copying
RUN ls -al

# Download Go modules
RUN go mod download

# Build the Go binary
RUN go build -o main .

# Prepare the dist directory with the built binary and necessary files
WORKDIR /dist
RUN cp /build/main .
RUN cp /build/*.properties .
RUN cp /build/*.json .
RUN ls /dist

# Step 2: Create the final image
FROM alpine:latest

# Install CA certificates in the final image
RUN apk add --no-cache ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary and other necessary files from the builder stage
COPY --from=builder /dist/main .
COPY --from=builder /dist/*.properties .
COPY --from=builder /dist/*.json .

# Expose the port if your application uses any
# EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/main"]