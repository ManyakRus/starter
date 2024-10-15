SERVICENAME=nikitin
SERVICEURL=gitlab.aescorp.ru/dsp_dev/claim/$(SERVICENAME)

FILEMAIN=./internal/v0/app/main.go
FILEAPP=./bin/app_race

NEW_REPO=github.com/ManyakRus/starter


run:
	clear
	./make_version.sh
	go build -race -o $(FILEAPP) $(FILEMAIN)
	#	cd ./bin && \
	./bin/app_race
mod:
	clear
	./make_version.sh
	go get -u ./...
	go mod tidy -compat=1.22
	go mod vendor
	go fmt ./...
build:
	clear
	./make_version.sh
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
lines:
	clear
	./make_version.sh
	go_lines_count ./ ./docs/lines_count.txt 10
licenses:
	golicense -out-xlsx=./docs/licenses.xlsx $(FILEAPP)
gocyclo:
	golangci-lint run ./... --disable-all -E gocyclo -v
