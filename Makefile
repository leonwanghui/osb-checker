BASE_DIR := $(shell pwd)
BUILD_DIR := $(BASE_DIR)/build

all: build

build: prebuild mockbroker

prebuild:
	mkdir -p $(BUILD_DIR)

mockbroker:
	go build -ldflags '-w -s' -o $(BUILD_DIR)/bin/osbchecker.mockbroker \
		github.com/openservicebrokerapi/osb-checker/mockbroker

autogenerated: precheck validate generate

precheck:
	wget https://raw.githubusercontent.com/openservicebrokerapi/servicebroker/master/swagger.yaml \
		-O $(BASE_DIR)/autogenerated/swagger.yaml

validate:
	docker run --rm -v $(BASE_DIR)/autogenerated:/local openapitools/openapi-generator-cli validate \
		-i /local/swagger.yaml

generate:
	cp $(BASE_DIR)/autogenerated/go-server/go/api_*.go $(BASE_DIR)/autogenerated/go-server

	docker run --rm -v $(BASE_DIR)/autogenerated:/local openapitools/openapi-generator-cli generate \
		-i /local/swagger.yaml \
		-g go \
		-o /local/go-client
	docker run --rm -v $(BASE_DIR)/autogenerated:/local openapitools/openapi-generator-cli generate \
    	-i /local/swagger.yaml \
		-g go-server \
		-o /local/go-server \
		--additional-properties hideGenerationTimestamp=true

	mv $(BASE_DIR)/autogenerated/go-server/api_*.go $(BASE_DIR)/autogenerated/go-server/go/

	rm -f $(BASE_DIR)/autogenerated/models/*.go && \
		cp $(BASE_DIR)/autogenerated/go-server/go/model_*.go $(BASE_DIR)/autogenerated/models/ && \
		rm -f $(BASE_DIR)/autogenerated/go-server/go/model_*.go

goimports:
	goimports -w $(shell go list -f {{.Dir}} ./... |grep -v /vendor/)

clean:
	rm -rf $(BUILD_DIR)

.PHONY: mockbroker autogenerated goimports clean
