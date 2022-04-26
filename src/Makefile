# List of implements rpc
PROJECT_RPCS = health/health.proto \
	safeweb_admin/service_auth.proto \
	safeweb_admin/service_user.proto \
	safeweb_admin/service_stock.proto

# List of implements server
SERVERS = $(PROJECT_RPCS)

# Enable grpc gateway
GATEWAY = true
GATEWAY_OPTIONS = ProtoJSONMarshaler Redirect

# For build commend
TARGET = bin
TARGET_BIN = safe-server.$$(date '+%m.%d.%H.%M')
TARGET_BIN_WINDOW = safe-server.exe
GO_CMD_MAIN = cmd/main.go

PROJECT_PROTO_FOLDER = proto
PROJECT_RPC_FOLDER = rpc

SERVER_PACKAGE_NAME =
SERVER_OUT_FOLDER = server

####################  DOES NOT EDIT BELLOW  ############################
.PHONY = build generate all clean

GO_TOOLS = golang.org/x/tools/cmd/stringer \
	github.com/stretchr/testify/mock \
	github.com/vektra/mockery/.../ \
	github.com/google/wire/cmd/wire \
	github.com/golang/protobuf/protoc-gen-go \
	github.com/gogo/protobuf/protoc-gen-gogo \
	github.com/gogo/protobuf/protoc-gen-gofast \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc \
	github.com/envoyproxy/protoc-gen-validate \
	github.com/mwitkow/go-proto-validators/protoc-gen-govalidators \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

$(GO_TOOLS):
	GOSUMDB=off go get $@

install-go-tools: $(GO_TOOLS)
	@echo \# installed go tools

# support fresh install on osx, not sure if it can't run properly
install-osx: $(GO_TOOLS)
	@brew install gnu-sed jq
	@export PATH="/usr/local/opt/gnu-sed/libexec/gnubin:$PATH"
	@brew install protobuf

# support fresh install on linux, not sure if it can't run properly
PROTOC_LINUX_VERSION = 3.11.4
PROTOC_LINUX_ZIP = protoc-$(PROTOC_LINUX_VERSION)-linux-x86_64.zip

install-linux:
	$(GO_TOOLS)
	@apt-get install -y jq
	@curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_LINUX_VERSION)/$(PROTOC_LINUX_ZIP)
	@sudo unzip -o $(PROTOC_LINUX_ZIP) -d /usr/local bin/protoc
	@sudo unzip -o $(PROTOC_LINUX_ZIP) -d /usr/local 'include/*'
	@rm -f $(PROTOC_LINUX_ZIP)

# support fresh install on window64, not sure if it can't run properly
PROTOC_WINDOW64_VERSION = 3.11.4
PROTOC_WINDOW64_ZIP = protoc-$(PROTOC_WINDOW64_VERSION)-win64.zip

install-window64: $(GO_TOOLS)
	@curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_WINDOW64_VERSION)/$(PROTOC_WINDOW64_ZIP)
	@sudo unzip $(PROTOC_WINDOW64_ZIP) -d /c/tools/protoc64
	@rm -f $(PROTOC_WINDOW64_ZIP)

.ONESHELL:
prepare-rpc:
	@mkdir -p $(PROJECT_RPC_FOLDER)

generate-rpc: prepare-rpc
	@echo \# generating rpc....
	@./swagger/_scripts/generate go
	@# Workaround for https://github.com/grpc-ecosystem/grpc-gateway/issues/229.
	@find ${PROJECT_RPC_FOLDER} -name "*.pb.gw.go" | xargs sed -i.bak "s/empty\./types\./g"
	@find ${PROJECT_RPC_FOLDER} -name "*.pb.gw.go" | xargs sed -i.bak "s/wrappers\./types\./g"
	@find ${PROJECT_RPC_FOLDER} -name "*.pb.gw.go" | xargs sed -i.bak "s/empty\ \"github.com\/golang\/protobuf\/ptypes\/empty\"/\"github.com\/gogo\/protobuf\/types\"/g"
	@find ${PROJECT_RPC_FOLDER} -name "*.pb.gw.go" | xargs sed -i.bak "s/wrappers\ \"github.com\/golang\/protobuf\/ptypes\/wrappers\"/\"github.com\/gogo\/protobuf\/types\"/g"
	@find ${PROJECT_RPC_FOLDER} -name "*.pb.gw.go.bak" -exec rm -f {} \;
	@find ${PROJECT_RPC_FOLDER} -name "*_grpc.pb.go" | xargs sed -i.bak "s/empty\./types\./g"
	@find ${PROJECT_RPC_FOLDER} -name "*_grpc.pb.go" | xargs sed -i.bak "s/wrappers\./types\./g"
	@find ${PROJECT_RPC_FOLDER} -name "*_grpc.pb.go" | xargs sed -i.bak "s/empty\ \"github.com\/golang\/protobuf\/ptypes\/empty\"/\"github.com\/gogo\/protobuf\/types\"/g"
	@find ${PROJECT_RPC_FOLDER} -name "*_grpc.pb.go" | xargs sed -i.bak "s/wrappers\ \"github.com\/golang\/protobuf\/ptypes\/wrappers\"/\"github.com\/gogo\/protobuf\/types\"/g"
	@find ${PROJECT_RPC_FOLDER} -name "*_grpc.pb.go.bak" -exec rm -f {} \;

generate-swagger:
	@echo \# generating swagger....
	@./swagger/_scripts/generate docs

prepare:
	@mkdir -p $(SERVER_OUT_FOLDER)/$(SERVER_PACKAGE_NAME)

photon-server:
	@echo \# generating server....
	@protoc -I $(PROJECT_PROTO_FOLDER) \
		-I .third_party/github.com/grpc-ecosystem/grpc-gateway/ \
		-I .third_party/github.com/gogo/googleapis \
		--descriptor_set_out=safeweb.protoset \
		--include_imports \
		--rpc-server_out=gateway=$(GATEWAY),gateway_options="$(GATEWAY_OPTIONS)":$(SERVER_OUT_FOLDER) \
	$(SERVERS)
		@find ${SERVER_OUT_FOLDER} -name "rpc-server.pb.go" | xargs sed -i.bak "s/runtime\.NewServeMux./runtime\.NewServeMux(runtime\.WithProtoErrorHandler(CustomHTTPError)\,/g"
		@find ${SERVER_OUT_FOLDER} -name "rpc-server.pb.go" | xargs sed -i.bak "s/mapstructure:\"http\" yaml:\"grpc\"/mapstructure:\"http\" yaml:\"http\"/g"
		@find ${SERVER_OUT_FOLDER} -name "rpc-server.pb.go" | xargs sed -i.bak "s/mapstructure:\"gprc\"/mapstructure:\"grpc\"/g"
		@find ${SERVER_OUT_FOLDER} -name "rpc-server.pb.go.bak" -exec rm -f {} \;

generate-server: prepare photon-server

generate: generate-rpc generate-swagger generate-server
	@echo \# generate finished

build:
	@go build -o $(TARGET)/$(TARGET_BIN_WINDOW) $(GO_CMD_MAIN)
	@echo \# build Windows server finished

live:
	@./$(TARGET)/$(TARGET_BIN_WINDOW) server

build-live: build live

run: generate build-live

build-linux:
	@GOOS=linux GOARCH=amd64 go build -o $(TARGET)/$(TARGET_BIN) $(GO_CMD_MAIN)

build-server: generate-rpc generate-server build-linux
	@echo \# build Linux server finished

migrate:
	@echo \# make migrate name="$(name)"
	@./$(TARGET)/$(TARGET_BIN_WINDOW) migrate create $(name)

migrate-up:
	@./$(TARGET)/$(TARGET_BIN_WINDOW) migrate up

migrate-down-1:
	@./$(TARGET)/$(TARGET_BIN_WINDOW) migrate down 1

clean:
	@rm -rf ./$(SERVER_OUT_FOLDER)/$(SERVER_PACKAGE_NAME)/*.pb.go
	@rm -rf ./$(PROJECT_RPC_FOLDER)
	@rm -rf ./$(TARGET)
