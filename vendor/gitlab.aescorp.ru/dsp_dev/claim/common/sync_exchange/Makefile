## test
test:
	clear
	go fmt ./...
	#go test -vet=all -race -count 1 -test.shuffle on -coverprofile cover.out -covermode atomic ./... ./data_packer/... ./liveness/... ./sync_confirm/... ./sync_global/... ./sync_types/...
	go test -vet=all -race -count 1 -test.shuffle on -coverprofile cover.out -covermode atomic ./...
	#go test -vet=all -race -count 1 -test.shuffle on -coverprofile cover.out -covermode atomic ./data_packer/...
	#go test -vet=all -race -count 1 -test.shuffle on -coverprofile cover.out -covermode atomic ./liveness/...
	#go test -vet=all -race -count 1 -test.shuffle on -coverprofile cover.out -covermode atomic ./sync_confirm/...
	#go test -vet=all -race -count 1 -test.shuffle on -coverprofile cover.out -covermode atomic ./sync_global/...
	#go test -vet=all -race -count 1 -test.shuffle on -coverprofile cover.out -covermode atomic ./sync_types/...
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
	# wget http://support.aescorp.ru/.golangci.yml
	golangci-lint run
	#golangci-lint run ./data_packer/... ./sync_confirm/... ./sync_global/... ./sync_types/...
	#gocyclo -over 15 ./data_packer/
	#gocyclo -over 15 ./sync_confirm/
	#gocyclo -over 15 ./sync_global/
	#gocyclo -over 15 ./sync_types/
	#gocritic check ./data_packer/... ./sync_confirm/... ./sync_global/... ./sync_types/...
	#staticcheck ./data_packer/... ./sync_confirm/... ./sync_global/... ./sync_types/...

## help
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
