swagger-install:
	go get -u github.com/swaggo/swag/cmd/swag@v1.16.3
	go install github.com/swaggo/swag/cmd/swag@v1.16.3

swag:
	swag init -g api/api.go -o api/docs

run:
	go run cmd/main.go