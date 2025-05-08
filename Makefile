PROTO_PATH=internal/proto
PROTO_FILE_PATH=internal/proto/subpub.proto
BUILD_PATH=bin

.PHONY: test clean

all: build-server build-client

build-server:
	@echo Building server...
	@go build -o $(BUILD_PATH)/server ./cmd/server

build-client:
	@echo Building client...
	@go build -o $(BUILD_PATH)/client ./test/client

proto_gen:
	@echo Successful generated
	@protoc -I internal/proto/ internal/proto/subpub.proto \
	--go_out=internal/proto/gen/ \
	--go_opt=paths=source_relative \
	--go-grpc_out=internal/proto/gen/ \
	--go-grpc_opt=paths=source_relative

test:
	go test ./pkg/subpub

clean:
	@echo Done!
	@rm -rf $(BUILD_PATH)

