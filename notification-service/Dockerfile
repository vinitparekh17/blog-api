FROM golang:1.21.5-alpine3.18 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./docs /app/docs
COPY ./ .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -o notification-api 

FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=builder /app/notification-api notification-api
EXPOSE 8080
CMD [ "./notification-api" ]