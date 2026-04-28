FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
COPY data-model/go.mod data-model/go.sum ./data-model/
COPY swagger/go.mod swagger/go.sum ./swagger/
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/run_service ./cmd/run_service

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /bin/run_service /app/run_service

EXPOSE 8380

ENTRYPOINT ["/app/run_service"]
