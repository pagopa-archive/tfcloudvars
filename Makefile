.DEFAULT_GOAL := build

BIN_FILE=tfcloudvars

build:
	@go build -o "${BIN_FILE}"
clean:
	go clean
	rm --force "cp.out"
	rm --force nohup.out
test:
	go test ./...
