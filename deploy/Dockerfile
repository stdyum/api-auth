FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN rm -f go.work go.work.sum
RUN go build -o app github.com/stdyum/api-auth/cmd/auth

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/scheme ./scheme
CMD ["./app"]