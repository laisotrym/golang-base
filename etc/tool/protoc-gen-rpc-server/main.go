package main

import (
	"flag"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"jupiter.app/tool/protoc-gen-rpc-server/internal_gen"
)

func main() {

	var (
		importPath         string
		gatewayEnabled     bool
		gatewayOptionsText string
	)

	flag.StringVar(&importPath, "import_path", "rpc.safeweb.app/server", "Output of gRPC server")
	flag.BoolVar(&gatewayEnabled, "gateway", false, "Enable grpc-gateway support")
	flag.StringVar(&gatewayOptionsText, "gateway_options", "DefaultMarshaler", "List of GRPC-Gateway server options")

	protogen.Options{
		ParamFunc: flag.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		g := plugin.NewGeneratedFile("rpc-server.pb.go", protogen.GoImportPath(importPath))
		gatewayOptions := strings.Split(gatewayOptionsText, " ")
		return internal_gen.GenerateFile(plugin, protogen.GoImportPath(importPath), gatewayEnabled, gatewayOptions, g)
	})
}
