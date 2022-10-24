# syntax=docker/dockerfile:1

## Build
FROM golang:alpine AS build
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /avito-internship

## Deploy
FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /
COPY --from=build /avito-internship /avito-internship
EXPOSE 8080
ENTRYPOINT ["/avito-internship"]
