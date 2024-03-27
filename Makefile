# Set env var to "-short" to only run short tests. Use GO_TEST_SHORT env var from pre-commit
_GO_TEST_SHORT?=

mocks: clean_mocks
	go install github.com/golang/mock/mockgen@v1.6.0
	go generate ./...

clean_mocks:
	find . -name mocks -type d  -exec rm -r {} +

tests:
	go test ./... -v
