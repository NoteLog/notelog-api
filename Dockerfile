FROM golang:latest as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v

FROM alpine:3
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/notelog-api /notelog-api

CMD ["./notelog-api"]
