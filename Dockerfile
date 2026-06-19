# syntax=docker/dockerfile:1

ARG GO_VERSION=1.25

FROM golang:${GO_VERSION}-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/suprise ./cmd/server

FROM scratch

COPY --from=build /out/suprise /suprise

EXPOSE 8080

ENTRYPOINT ["/suprise"]
