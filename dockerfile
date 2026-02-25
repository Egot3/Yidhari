FROM golang:1.26.0-alpine AS builder
RUN apk add --no-cache git
WORKDIR /Yidhari
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app .

FROM scratch
COPY --from=builder /app /app
ENV PORT=380 URL=amqp://guest:guest@host.docker.internal OWN_PORT=9130 REPORTING_PORT=8250
EXPOSE 9130 8250 380
ENTRYPOINT ["/app"]