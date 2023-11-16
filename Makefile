compose_up:
	docker compose up

compose_down:
	docker compose down

test_sync_account_limit:
	curl 127.0.0.1:8080/flush_redis
	go run test/test_account_sync.go
test_async_account_limit:
	curl 127.0.0.1:8080/flush_redis
	go run test/test_account_async.go

test_sync_endpoint_limit:
	curl 127.0.0.1:8080/flush_redis
	go run test/test_endpoint_sync.go
test_async_endpoint_limit:
	curl 127.0.0.1:8080/flush_redis
	go run test/test_endpoint_async.go
