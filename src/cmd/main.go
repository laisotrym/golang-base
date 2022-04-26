package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	safeweb_lib_interceptor "safeweb.app/service/safeweb_lib/interceptor"

	"safeweb.app/service/safeweb_admin/auth"
	"safeweb.app/service/safeweb_admin/core/stock"
	"safeweb.app/service/safeweb_admin/core/user"

	"github.com/gogo/status"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"

	"safeweb.app/database/migrate"

	"safeweb.app/config"
	"safeweb.app/server"
	"safeweb.app/service/health"
)

var (
	cmdRoot = &cobra.Command{
		Use: "backend-api",
	}
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	conf, err := config.Load()
	if err != nil {
		return err
	}
	logger, err := conf.Log.Build()
	if err != nil {
		return err
	}
	fmt.Printf("Set timezone: %s\n", conf.TimeZone)
	config.SetTimeZone(&conf.TimeZone)

	cmdList := make([]*cobra.Command, 1+len(conf.App))
	cmdIdx := 0
	for k, v := range conf.App {
		if k != config.AppAliasDefault {
			cmd := &cobra.Command{
				Use:   fmt.Sprintf("%s", strings.Replace(strings.ToLower(k.ToString()), "_", "-", -1)),
				Short: fmt.Sprintf("Command for %s service", strings.ToUpper(strings.Replace(k.ToString(), "service_", "", -1))),
			}
			cmd.AddCommand(migrate.Command(v.MigrationsFolder, v.MySQL.String()))
			cmdList[cmdIdx] = cmd
			cmdIdx++
		}
	}

	cmdList[cmdIdx] = &cobra.Command{
		Use: "server",
		Run: func(cmd *cobra.Command, args []string) {
			grpc_zap.ReplaceGrpcLoggerV2(logger)

			opts := []grpc_recovery.Option{
				grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
					logger.Error("recovery err: "+string(debug.Stack()), zap.Any("error", p))
					return status.Error(codes.Unknown, "panic triggered: ")
				}),
			}

			jwtManager := auth.NewJWTManager(config.SecretKey, config.TokenDuration)
			authInterceptor := safeweb_lib_interceptor.NewInterceptor(jwtManager, config.AccessibleRoles())
			s := server.NewServer(conf.Server,
				grpc_middleware.WithUnaryServerChain(
					grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
					grpc_recovery.UnaryServerInterceptor(opts...),
					authInterceptor.UnaryAuthInterceptor(),
				),
				grpc_middleware.WithStreamServerChain(
					grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
					grpc_recovery.StreamServerInterceptor(opts...),
					authInterceptor.StreamAuthInterceptor(),
				),
			)

			heart := health.NewServer()
			heart.Beat("server")
			var serversBackend []interface{}
			if v, ok := conf.App[config.AppAliasBackend]; ok {
				dbConn, _, err := config.InitConn(v, false, true)
				if err != nil {
					logger.Fatal("Init connection error", zap.String("App", config.AppAliasBackend.ToString()), zap.Any("error", err), zap.String("DSN", v.MySQL.DSN()))
				}

				serversBackend = []interface{}{
					user.InitializeServer(dbConn),
					auth.InitializeServer(dbConn, jwtManager),
					stock.InitializeServer(dbConn),
				}
			}

			// Init connection
			servers := make([]interface{}, len(serversBackend)+1)
			serverIdx := 0
			servers[serverIdx] = heart
			serverIdx++
			servers = addServers(&serverIdx, serversBackend, servers)

			// Register your rpc server here
			// You may register multiple server
			//
			if err := s.Register(servers...); err != nil {
				logger.Fatal("Error register servers", zap.Any("error", err))
			}

			if err := s.WebServe(conf.Server.HTTP.String(), conf.Server.GRPC.String(), conf.Webroot); err != nil {
				logger.Fatal("Error start server", zap.Any("error", err))
			}
		},
	}
	cmdIdx++
	cmdList[cmdIdx] = migrate.Command(conf.MigrationsFolder, conf.MySQL.String())

	cmdRoot.AddCommand(cmdList...)
	return cmdRoot.Execute()
}

func addServers(idx *int, servers []interface{}, output []interface{}) []interface{} {
	for _, v := range servers {
		output[*idx] = v
		*idx++
	}

	return output
}
