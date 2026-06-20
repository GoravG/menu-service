from golang:1.26.4-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o main .

FROM alpine:latest

COPY --from=builder /app/main /app/main
RUN chmod +x /app/main

EXPOSE 8080

CMD ["/app/main"]