FROM golang:1.26-alpine@sha256:2389ebfa5b7f43eeafbd6be0c3700cc46690ef842ad962f6c5bd6be49ed82039 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/self/main.go ./
RUN CGO_ENABLED=0 go build -o smoke-app .

FROM scratch
COPY --from=builder /app/smoke-app /smoke-app
ENTRYPOINT ["/smoke-app"]
