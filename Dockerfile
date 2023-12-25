FROM golang:1.21.5-bookworm AS builder
WORKDIR /server
COPY go.* .
RUN go mod download
COPY . .
RUN go build -o app ./cmd/server

FROM gcr.io/distroless/base-debian12:debug
WORKDIR /app
COPY --from=builder /server/app .
RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
EXPOSE 8080
CMD [ "./app" ]
