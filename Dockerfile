FROM golang:1.27.1-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server .

FROM debian:trixie

WORKDIR /srv

COPY --from=builder /app/server .
RUN apt-get update && apt-get install ca-certificates
RUN mkdir -p ./public

CMD [ "./server" ]