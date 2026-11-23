FROM golang:1.26-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd

FROM alpine:3.21
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=build /app/server ./server
COPY .env.example ./.env
EXPOSE 9000
CMD ["./server"]
