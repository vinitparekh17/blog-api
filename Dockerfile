FROM golang:1.21.5-alpine3.18 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/root/.cache/go-build go build -o main cmd/api/blog/main.go

FROM alpine:3.14.2
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main 1"]