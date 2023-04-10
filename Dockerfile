FROM golang:1.20.3-alpine3.16 AS builder

WORKDIR /app

# Copy files
COPY go.mod ./
COPY magicpacket ./magicpacket
COPY main.go ./

# Add group and user
RUN addgroup -S wol-api && adduser -S wol-api -G wol-api

# Build binary
RUN CGO_ENABLED=0 go build -o wol-api


FROM scratch

ENV LISTEN_ADDRESS="0.0.0.0"

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/wol-api /usr/bin/wol-api

USER wol-api

EXPOSE 9099

ENTRYPOINT ["/usr/bin/wol-api"]
