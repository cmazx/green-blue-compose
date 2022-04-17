FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY ./main.go ./main.go
RUN GOOS=linux go build -o /app/app_bin /app/main.go
CMD ["./app_bin"]
