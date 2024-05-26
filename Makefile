run: 
	@go run .

docker-build:
	 @docker build -t gomysql_img .

docker-run:
	@docker run --rm --name go_mysql -p 3306:3306 gomysql_img

test:
	@go test -v ./...