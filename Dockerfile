FROM golang:1.22-alpine as build

WORKDIR /app

COPY . .

RUN go mod download

RUN go get -u github.com/swaggo/swag/cmd/swag@v1.16.3
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3
RUN swag init -g api/api.go -o api/docs

RUN go build -o medods cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/medods .
COPY --from=build /app/migrations /root/migrations

EXPOSE 8095

CMD ["./medods"]