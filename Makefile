PROJECTNAME=$(shell basename "$(PWD)")
GOLANGCI := $(GOPATH)/bin/golangci-lint

.PHONY: help lint test
all: help
help: Makefile
	@echo
	@echo " Choose a make command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

$(GOLANGCI):
	if [ ! -f ./bin/golangci-lint ]; then \
		wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest; \
	fi;

lint: $(GOLANGCI)
	./bin/golangci-lint run ./... --timeout 5m0s

test:
	go test ./... 

test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -html=coverage.out


docker:
	docker-compose -f ./docker-compose-chains.yml up -V

## license: Adds license header to missing files.
license:
	@echo "  >  \033[32mAdding license headers...\033[0m "
	GO111MODULE=off go get -u github.com/google/addlicense
	addlicense -c "ChainSafe Systems" -f ./copyright.txt -y 2020 .

## license-check: Checks for missing license headers
license-check:
	@echo "  >  \033[Checking for license headers...\033[0m "
	GO111MODULE=off go get -u github.com/google/addlicense
	addlicense -check -c "ChainSafe Systems" -f ./copyright.txt -y 2020 .

rebuild-contracts:
	rm -rf bindings/ solidity/
	TARGET=build ./scripts/setup_contracts.sh

clean:
	rm -rf build/ solidity/

start-elections:
	docker-compose -f ./docker-compose-elections.yml up -V


build:
	@echo "  >  \033[32mBuilding binary...\033[0m "
	GOARCH=amd64 go build -o build/chainbridge-celo

install:
	@echo "  >  \033[32mInstalling bridge...\033[0m "
	go install


genmocks:
	mockgen -destination=./chain/listener/mock/listener.go -source=./chain/listener/listener.go
	mockgen -destination=./chain/listener/mock/bindings.go -source=./chain/listener/bindings.go -package=mock_listener
	mockgen -destination=./chain/writer/mock/writer.go -source=./chain/writer/writer.go
	mockgen -destination=./chain/mock/chain.go -source=./chain/chain.go
	mockgen -destination=./chain/client/mock/client.go -source=./chain/client/client.go



