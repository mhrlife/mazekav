test:
	go test ./... -v
test-integration:
	TEST_TYPE=integration go test ./... -v