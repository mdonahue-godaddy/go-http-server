FROM alpine:latest

WORKDIR /opt/http/server

COPY . .

CMD ["./go-http-server"]
