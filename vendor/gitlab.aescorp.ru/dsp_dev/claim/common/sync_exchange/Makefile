## test
test:
	clear
	go fmt ./...
	go test -vet=all -race -count 1 -test.shuffle on -coverprofile cover.out -covermode atomic ./... ./data_packer/... ./sync_confirm/... ./sync_global/... ./sync_types/...
	go tool cover -func=cover.out

## mod
mod:
	clear
	go mod tidy -compat=1.22
	go mod vendor
	go fmt ./...

## commit
commit:
	git commit -m "Update version"
	git push origin master
	git tag "v0.0.31"
	git push origin master --tags

## lint
lint:
	clear
	go fmt ./...
	golangci-lint run ./data_packer/... ./sync_confirm/... ./sync_global/... ./sync_types/...
	gocyclo -over 15 ./data_packer/
	gocyclo -over 15 ./sync_confirm/
	gocyclo -over 15 ./sync_global/
	gocyclo -over 15 ./sync_types/
	gocritic check ./data_packer/... ./sync_confirm/... ./sync_global/... ./sync_types/...
	staticcheck ./data_packer/... ./sync_confirm/... ./sync_global/... ./sync_types/...

## help
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
