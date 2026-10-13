FROM golang:1.26-alpine@sha256:3082400e369fa24d5fc60bca20edab3f6d604e0c5a690ec66b295eff4dd87ade AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/self/main.go ./
RUN CGO_ENABLED=0 go build -o smoke-app .

FROM scratch
COPY --from=builder /app/smoke-app /smoke-app
ENTRYPOINT ["/smoke-app"]
