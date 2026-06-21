main_package_path = ./
binary_name = main
binary_dir = tmp/bin
binary_path = $(binary_dir)/$(binary_name)

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: run quality control checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## test: run all tests
.PHONY: test
test:
	@set -a && [ -f .env ] && . ./.env; set +a && go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## tidy: tidy modfiles and format .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fix ./...
	go fmt ./...

## build: build the application
.PHONY: build
build:
	@mkdir -p $(binary_dir)
	go build -o=$(binary_path) ${main_package_path}

## run: run the application
.PHONY: run
run: build
	@set -a && [ -f .env ] && . ./.env; set +a && ./$(binary_path)

## live: run the application with reloading on file changes
.PHONY: live
live:
	@set -a && [ -f .env ] && . ./.env; set +a && go tool air \
		--build.cmd "make build" --build.bin "$(binary_path)" --build.delay "100" \
		--build.include_ext "go" \
		--misc.clean_on_exit "true"
