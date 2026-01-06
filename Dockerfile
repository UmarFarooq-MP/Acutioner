FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o auctioner

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/auctioner /app/auctioner
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/auctioner"]
