FROM golang:1.25 AS builder
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download
COPY . .
ARG SERVICE
RUN test -n "$SERVICE" || (echo "SERVICE build arg is required" && exit 1)
RUN CGO_ENABLED=0 GOOS=linux \
    go build -o app ./cmd/${SERVICE}

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/app /app
ENTRYPOINT ["/app"]

