FROM golang:1.21 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o file-server ./cmd/app

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/file-server .
COPY /wait-for-pg.sh ./
RUN mkdir "files"
RUN apk --no-cache add postgresql-client && chmod +x wait-for-pg.sh
EXPOSE 8080
CMD ["./file-server"]
