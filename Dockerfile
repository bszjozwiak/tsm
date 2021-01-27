FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN go build .

FROM alpine:latest
COPY --from=builder /go/src/app/tsm .
CMD ["./tsm"]