# stage 0: compile go program
FROM golang:1.20
RUN mkdir -p /tmp/dr-gateway
WORKDIR /tmp/dr-gateway
ADD internal ./internal
ADD pkg ./pkg
ADD go.mod .
ADD go.sum .
RUN ls -l /tmp/dr-gateway && GOOS=linux go build -a -installsuffix cgo -o bin/dr-gateway-api internal/api-server/main.go

# stage 1: build image for the api-server
FROM alpine:latest as api-server
WORKDIR /root
EXPOSE 8080
COPY --from=0 /tmp/dr-gateway/bin/dr-gateway-api .

## entrypoint in shell form so that we can use $PORT environment variable
ENTRYPOINT ["./dr-gateway-api"]