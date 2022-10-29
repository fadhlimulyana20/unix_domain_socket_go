FROM golang:alpine as builder
LABEL maintainer="Fadhli Mulyana <baba.fadhli@gmail.com>"
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build main.go
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
ENTRYPOINT [ "./main" ]