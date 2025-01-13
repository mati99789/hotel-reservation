build:
	@go build -o ./bin/api
run: build
	@./bin/api
seed:
	@go run scripts/seed.go
test:
	@go test -v ./... -count=1

docker:
	echo "Building a docker file"
	@docker build -t api .
	echo "Running API inside Docker container"
	@docker run -p 3000:3000 api
