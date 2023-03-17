## run.test
run.test:
	clear
	go fmt ./...
	go test -coverprofile cover.out -covermode atomic ./ ./object_storage/... ./object_view/...
	go tool cover -func=cover.out

## mod
mod:
	clear
	go mod tidy -compat=1.17
	go mod vendor
	go fmt ./...

## lint
lint:
	clear
	go fmt ./...
	golangci-lint run ./ ./object_storage/... ./object_view/...
	gocyclo -over 15 ./
	gocyclo -over 15 ./object_storage/
	gocyclo -over 15 ./object_view/
	gocritic check ./ ./object_storage/... ./object_view/...
	staticcheck ./ ./object_storage/... ./object_view/....

## help
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
