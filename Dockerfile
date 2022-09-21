FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY ./cmd .
COPY ./pkg .

# Build the application
RUN go build -o go-http-server ./cmd/go-http-server/main.go

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/go-http-server .
COPY .config/go-http-server.json .

EXPOSE 8081
EXPOSE 8082

# Build image
FROM alpine:latest

COPY --from=builder /dist/go-http-server* /

ENTRYPOINT ["/go-http-server"]
