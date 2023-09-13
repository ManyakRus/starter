SERVICENAME=nikitin
SERVICEURL=gitlab.aescorp.ru/dsp_dev/claim/$(SERVICENAME)

FILEMAIN=./internal/v0/app/main.go
FILEAPP=./bin/app_race

NEW_REPO=github.com/ManyakRus/starter


run:
	clear
	go build -race -o $(FILEAPP) $(FILEMAIN)
	#	cd ./bin && \
	./bin/app_race
mod:
	clear
	go get -u ./...
	go mod tidy -compat=1.20
	go mod vendor
	go fmt ./...
build:
	clear
	go build -race -o $(FILEAPP) $(FILEMAIN)
	cd ./cmd && \
	./VersionToFile.py
lint:
	clear
	go fmt ./...
	golangci-lint run ./...
	gocyclo -over 10 ./
	gocritic check ./...
	staticcheck ./...
test.run:
	clear
	go fmt ./...
	go test -coverprofile cover.out ./...
	go tool cover -func=cover.out
newrepo:
	sed -i 's+$(SERVICEURL)+$(NEW_REPO)+g' go.mod
	find -name *.go -not -path "*/vendor/*"|xargs sed -i 's+$(SERVICEURL)+$(NEW_REPO)+g'
graph:
	clear
	image_packages ./ docs/packages.graphml
conn:
	clear
	image_connections ./ docs/connections.graphml $(SERVICENAME)
