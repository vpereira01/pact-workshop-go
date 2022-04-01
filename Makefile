TEST?=./...

include ./make/config.mk

#install:
#	@echo "--- Installing Pact CLI dependencies"
#	curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash

run-consumer:
	@go run consumer/client/cmd/main.go

run-provider:
	@go run provider/cmd/usersvc/main.go

unit:
	@echo "--- 🔨Running Unit tests "
	cd consumer/client; go test -v .

pact-consumer:
	@echo "--- 🔨Running Consumer Pact tests "
	cd consumer/client; go test -v -tags=pact -count=1 -timeout=1m . -run=Pact

pact-provider:
	@echo "--- 🔨Running Provider Pact tests "
	cd provider; go test -v -tags=pact -count=1 -timeout=1m . -run=Pact

#.PHONY: install unit consumer  run-provider run-consumer
.PHONY: unit consumer  run-provider run-consumer