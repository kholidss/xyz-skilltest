BINARY_NAME=xyz-skilltest
build:
	@go build -o bin/${BINARY_NAME} main.go

run-http:
	@./bin/${BINARY_NAME} http

run-rabbit:
	@./bin/${BINARY_NAME} rabbit
	
install:
	@echo "Installing dependencies...."
	@rm -rf vendor
	@rm -f Gopkg.lock
	@rm -f glide.lock
	@go mod tidy && go mod download && go mod vendor

test:
	@go test $$(go list ./... | grep -v /vendor/) -cover

test-cover:
	@go test $$(go list ./... | grep -v /vendor/) -coverprofile=cover.out && go tool cover -html=cover.out ; rm -f cover.out

coverage:
	@go test -covermode=count -coverprofile=count.out fmt; rm -f count.out

start-http:
	@echo "Starting HTTP Service...."
	@go run main.go http

start-rabbit:
	@echo "Starting RabbitMQ Listener...."
	@go run main.go rabbit

migration:
	@go run main.go db:migrate

migration-create:
	@go run main.go db:migrate create $(name) sql

migration-status:
	@go run main.go db:migrate status

migration-up:
	@echo "Starting database migration: UP...."
	@go run main.go db:migrate up

migration-down:
	@echo "Starting database migration: DOWN...."
	@go run main.go db:migrate down