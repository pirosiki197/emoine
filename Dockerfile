FROM golang:1.21.1-alpine AS builder
WORKDIR /server
COPY go.* .
RUN go mod download
COPY . .
RUN go build -o app ./cmd/server

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /server/app .
RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
EXPOSE 8080
CMD [ "./app" ]
