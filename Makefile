PROTO_PATH=internal/proto
PROTO_FILE_PATH=internal/proto/subpub.proto
PROTO_GEN_PATH=internal/proto/gen
BUILD_PATH=bin

.PHONY: all build-server build-client test detailed_test clean gen

all: build-server build-client

build-server:
	@echo Building server...
	@go build -o $(BUILD_PATH)/server ./cmd/server

build-client:
	@echo Building client...
	@go build -o $(BUILD_PATH)/client ./test/client

gen:
	@echo Successful generated
	@protoc -I $(PROTO_PATH) $(PROTO_FILE_PATH) \
	--go_out=$(PROTO_GEN_PATH) \
	--go_opt=paths=source_relative \
	--go-grpc_out=$(PROTO_GEN_PATH) \
	--go-grpc_opt=paths=source_relative

test:
	@go test ./pkg/subpub

detailed_test:
	@go test -v ./pkg/subpub

clean:
	@rm -rf $(BUILD_PATH)
	@echo Done!


