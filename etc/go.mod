module jupiter.app

go 1.14

require (
	github.com/dave/jennifer v1.4.0
	github.com/duythinht/dbml-go v0.0.0-20200518103229-5472f2db3240
	github.com/golang/protobuf v1.4.2
	github.com/jinzhu/gorm v1.9.14
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	google.golang.org/genproto v0.0.0-20200519141106-08726f379972
	google.golang.org/protobuf v1.23.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace rpc.tekoapis.com => ./third_party/rpc.tekoapis.com
