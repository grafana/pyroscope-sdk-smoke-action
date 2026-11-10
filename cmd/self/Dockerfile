FROM golang:1.26-alpine@sha256:30b1e546a9c7fa3b747c2bdba3a0c7bdf3a02ebb4c5b8405d2f948d7e8aac4ec AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/self/main.go ./
RUN CGO_ENABLED=0 go build -o smoke-app .

FROM scratch
COPY --from=builder /app/smoke-app /smoke-app
ENTRYPOINT ["/smoke-app"]
