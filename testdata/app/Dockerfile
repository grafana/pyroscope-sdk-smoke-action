FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY testdata/app/go.mod testdata/app/go.sum ./
RUN go mod download
COPY testdata/app/main.go ./
RUN CGO_ENABLED=0 go build -o smoke-app .

FROM scratch
COPY --from=builder /app/smoke-app /smoke-app
ENTRYPOINT ["/smoke-app"]
