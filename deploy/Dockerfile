FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN rm -f go.work go.work.sum
RUN go build -o app github.com/stdyum/api-auth/cmd/auth

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/scheme ./scheme

HEALTHCHECK --interval=1s --retries=120 CMD wget --no-verbose --tries=1 --output-document=/dev/null http://0.0.0.0:$(echo "${PORT_HTTP}" | sed -e 's/^"//' -e 's/"$//' -e "s/^'//" -e "s/'$//")/ping || exit 1
CMD ["./app"]