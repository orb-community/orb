package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ns1labs/orb/buildinfo"
	"github.com/ns1labs/orb/migrate"
	"github.com/ns1labs/orb/migrate/migration"
	"github.com/ns1labs/orb/migrate/postgres"
	"github.com/ns1labs/orb/pkg/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

const (
	svcName   = "migrate"
	envPrefix = "orb_migrate"
)

var log *zap.Logger

func init() {
	atomicLevel := zap.NewAtomicLevel()
	svcCfg := config.LoadBaseServiceConfig(envPrefix, "")

	switch strings.ToLower(svcCfg.LogLevel) {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stdout,
		atomicLevel,
	)

	log = zap.New(core, zap.AddCaller())
	log.Info("initialising logger")
}

func main() {
	ketoDbCfg := config.LoadPostgresConfig(fmt.Sprintf("%s_%s", envPrefix, postgres.DbKeto), postgres.DbKeto)
	usersDbCfg := config.LoadPostgresConfig(fmt.Sprintf("%s_%s", envPrefix, postgres.DbUsers), postgres.DbUsers)
	thingsDbCfg := config.LoadPostgresConfig(fmt.Sprintf("%s_%s", envPrefix, postgres.DbThings), postgres.DbThings)
	sinksDbCfg := config.LoadPostgresConfig(fmt.Sprintf("%s_%s", envPrefix, postgres.DBSinks), postgres.DBSinks)
	sinksEncryptionKey := config.LoadEncryptionKey(fmt.Sprintf("%s_%s", envPrefix, postgres.DBSinks))

	dbs := make(map[string]postgres.Database)

	dbs[postgres.DbKeto] = connectToDB(ketoDbCfg, true, log)
	dbs[postgres.DbUsers] = connectToDB(usersDbCfg, false, log)
	dbs[postgres.DbThings] = connectToDB(thingsDbCfg, false, log)

	sinksDB := connectToDB(sinksDbCfg, false, log)

	svc := migrate.New(
		log,
		dbs,
		migration.NewM1KetoPolicies(log, dbs),
		migration.NewM2SinksCredentials(log, sinksDB, sinksEncryptionKey),
	)

	rootCmd := &cobra.Command{
		Use: "orb-migrate",
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show migrate version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("orb-migrate %s\n", buildinfo.GetVersion())
		},
	}

	upCmd := &cobra.Command{
		Use:   "up",
		Short: "Execute migrations considering the existent schema version",
		Long:  "Execute migrations considering the existent schema version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := svc.Up(); err != nil {
				log.Error("error executing migration up", zap.Error(err))
				os.Exit(1)
			}
		},
	}

	downCmd := &cobra.Command{
		Use:   "down",
		Short: "Rollback migrations considering the existent schema version",
		Long:  "Rollback migrations considering the existent schema version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := svc.Down(); err != nil {
				log.Error("error executing migration down", zap.Error(err))
				os.Exit(1)
			}
		},
	}

	dropCmd := &cobra.Command{
		Use:   "drop",
		Short: "Rollback all migrations",
		Long:  "Rollback all migrations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := svc.Drop(); err != nil {
				log.Error("error executing migration drop", zap.Error(err))
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(dropCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Error("error on command exit", zap.Error(err))
		fmt.Printf("error on command exit: %s\n", err.Error())
		os.Exit(1)
	}
}

func connectToDB(cfg config.PostgresConfig, migrate bool, logger *zap.Logger) *sqlx.DB {
	db, err := postgres.Connect(cfg, migrate)
	if err != nil {
		logger.Error("Failed to connect to postgres", zap.Error(err))
		os.Exit(1)
	}
	return db
}
