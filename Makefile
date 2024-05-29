ENV_LOCAL_TEST=\
  TESTCONTAINERS_RYUK_DISABLED=true

run:
	go run cmd/main.go

test:
	$(ENV_LOCAL_TEST) \
      go test -v ./...

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down



