FROM golang:1.23-alpine@sha256:383395b794dffa5b53012a212365d40c8e37109a626ca30d6151c8348d380b5f AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/self/main.go ./
RUN CGO_ENABLED=0 go build -o smoke-app .

FROM scratch
COPY --from=builder /app/smoke-app /smoke-app
ENTRYPOINT ["/smoke-app"]
