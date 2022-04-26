// +build tools

package tools

import (
	_ "golang.org/x/tools/cmd/stringer"
	_ "github.com/stretchr/testify/mock"
	_ "github.com/vektra/mockery"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/gogo/protobuf/protoc-gen-gogo"
	_ "github.com/gogo/protobuf/protoc-gen-gofast"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/mwitkow/go-proto-validators/protoc-gen-govalidators"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
	_ "github.com/rakyll/statik"
)
