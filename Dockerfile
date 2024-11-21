FROM golang:alpine3.20 AS build
WORKDIR /app
COPY ./cmd /app/cmd
COPY ./src /app/src
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bridger ./cmd/bridger/main.go

FROM alpine:3.20.3
RUN apk --update add --no-cache openssh jq
COPY --from=build /bridger /usr/bin/bridger

