package migrate

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	migrateV4 "github.com/golang-migrate/migrate/v4"

	// import mysql
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	// import posgres
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// import file
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const versionTimeFormat = "20060102150405"

// Command return common migrate command for application
func Command(sourceURL string, databaseURL string) *cobra.Command {
	//Migration should always run on development mode
	logger, _ := zap.NewDevelopment()

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "database migration command",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "lift migration up to date",
		Run: func(cmd *cobra.Command, args []string) {
			m, err := migrateV4.New(sourceURL, databaseURL)

			if err != nil {
				logger.Fatal("Error create migration", zap.Error(err))
			}

			logger.Info("migration up")
			if err := m.Up(); err != nil && err != migrateV4.ErrNoChange {
				logger.Fatal(err.Error())
			}
		},
	}, &cobra.Command{
		Use:   "down",
		Short: "step down migration by N(int)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			m, err := migrateV4.New(sourceURL, databaseURL)

			if err != nil {
				logger.Fatal("Error create migration", zap.Error(err))
			}

			down, err := strconv.Atoi(args[0])

			if err != nil {
				logger.Fatal("rev should be a number", zap.Error(err))
			}

			logger.Info("migration down", zap.Int("down", -down))
			if err := m.Steps(-down); err != nil {
				logger.Fatal(err.Error())
			}
		},
	}, &cobra.Command{
		Use:   "force",
		Short: "Enforce dirty migration with verion (int)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			m, err := migrateV4.New(sourceURL, databaseURL)

			if err != nil {
				logger.Fatal("Error create migration", zap.Error(err))
			}

			ver, err := strconv.Atoi(args[0])

			if err != nil {
				logger.Fatal("rev should be a number", zap.Error(err))
			}

			logger.Info("force", zap.Int("ver", ver))

			if err := m.Force(ver); err != nil {
				logger.Fatal(err.Error())
			}
		},
	}, &cobra.Command{
		Use:  "create",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			folder := strings.ReplaceAll(sourceURL, "file://", "")
			now := time.Now()
			ver := now.Format(versionTimeFormat)
			name := strings.Join(args, "-")

			up := fmt.Sprintf("%s/%s_%s.up.sql", folder, ver, name)
			down := fmt.Sprintf("%s/%s_%s.down.sql", folder, ver, name)

			logger.Info("create migration", zap.String("name", name))
			logger.Info("up script", zap.String("up", up))
			logger.Info("down script", zap.String("down", up))

			if err := ioutil.WriteFile(up, []byte{}, 0644); err != nil {
				logger.Fatal("Create migration up error", zap.Error(err))
			}
			if err := ioutil.WriteFile(down, []byte{}, 0644); err != nil {
				logger.Fatal("Create migration down error", zap.Error(err))
			}
		},
	})
	return cmd
}
