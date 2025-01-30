#build stage
FROM golang:1.22 AS builder
ARG SERVICE_NAME
WORKDIR /usr/src/services

COPY go.mod go.sum ./

RUN go mod tidy && go mod download && go mod verify

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o bin/app cmd/$SERVICE_NAME/main.go

#final stage
FROM alpine:latest

WORKDIR /usr/bin/app/

COPY --from=builder /usr/src/services/bin/app .

CMD ["./app"]


