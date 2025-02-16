FROM golang:1.22.2-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./src ./src

RUN go build -o project src/main.go

FROM alpine:latest

COPY --from=builder ["/build/project", "/"]

ENTRYPOINT ["/project"]