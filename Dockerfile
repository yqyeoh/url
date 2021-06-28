FROM golang:1.16.3 AS builder

WORKDIR /build
RUN adduser -u 10001 -disabled-password -gecos '' app-user

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o url

FROM alpine:3.13 AS final

WORKDIR /app
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/url /app/url
COPY config /app/config

EXPOSE 3006

USER app-user
ENTRYPOINT ["/app/url"]
