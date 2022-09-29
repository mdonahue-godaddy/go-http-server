FROM golang:alpine AS builder

# Set necessary environmet variables needed for go build
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY ./cmd/go-http-server/main.go ./cmd/go-http-server/main.go
COPY ./ ./
COPY ./go.mod .
COPY ./go.sum .

# Build the application
RUN go build -o go-http-server ./cmd/go-http-server/main.go

# Move to /dist directory as the place for files for final container
WORKDIR /dist

# Copy binary from build to /dist folder
RUN cp /build/go-http-server ./go-http-server

# Copy app config json from src to /dist folder
COPY ./go-http-server.json ./go-http-server.json

# Build image
FROM alpine:latest

#EXPOSE 8081
#EXPOSE 8082

# Make app directory
RUN mkdir --verbose --parents /opt/go-http-server

# Move to /opt/go-http-server directory
WORKDIR /opt/go-http-server

# Copy files from builder /dist to WORKDIR  /opt/go-http-server
COPY --from=builder /dist/ .

ENTRYPOINT ["/opt/go-http-server/go-http-server"]
