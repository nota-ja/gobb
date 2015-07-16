package utils

import (
	"bitbucket.org/liamstask/goose/lib/goose"
	"fmt"
	"github.com/stevenleeg/gobb/config"
	"github.com/stevenleeg/gobb/models"
	"go/build"
	"path/filepath"
	"os"
)

var goose_conf *goose.DBConf

func generateGooseDbConf() *goose.DBConf {
	if goose_conf != nil {
		return goose_conf
	}

	pkg, _ := build.Import("github.com/stevenleeg/gobb/gobb", ".", build.FindOnly)
	db_username, _ := config.Config.GetString("database", "username")
	db_password, _ := config.Config.GetString("database", "password")
	db_database, _ := config.Config.GetString("database", "database")
	db_hostname, _ := config.Config.GetString("database", "hostname")
	db_port, _ := config.Config.GetString("database", "port")
	migrations_path := filepath.Join(pkg.SrcRoot, pkg.ImportPath, "../db/migrations")

	db_env_hostname, _ := config.Config.GetString("database", "env_hostname")
	db_env_port, _ := config.Config.GetString("database", "env_port")

	// Allow database information to come from environment variables
	if db_env_hostname != "" {
		db_hostname = os.Getenv(db_env_hostname)
	}
	if db_env_port != "" {
		db_port = os.Getenv(db_env_port)
	}

	if db_port == "" {
		db_port = "5432"
	}

	openstr := os.Getenv("DATABASE_URL")
	if openstr == "" {
		openstr = fmt.Sprintf("user=%s dbname=%s password=%s port=%s host=%s sslmode=disable", db_username, db_database, db_password, db_port, db_hostname)
	}

	goose_conf = &goose.DBConf{
		MigrationsDir: migrations_path,
		Env:           "development",
		Driver: goose.DBDriver{
			Name:    "postgres",
			OpenStr: openstr,
			Import:  "github.com/lib/pq",
			Dialect: &goose.PostgresDialect{},
		},
	}

	return goose_conf
}

func GetMigrationInfo() (latest_db_version int64, migrations []*goose.Migration, err error) {
	goose_conf := generateGooseDbConf()
	db := models.GetDbSession()

	latest_db_version, _ = goose.GetMostRecentDBVersion(goose_conf.MigrationsDir)
	current_db_version, _ := goose.EnsureDBVersion(goose_conf, db.Db)
	migrations, _ = goose.CollectMigrations(goose_conf.MigrationsDir, current_db_version, latest_db_version)

	return latest_db_version, migrations, err
}

func RunMigrations(version int64) error {
	goose_conf := generateGooseDbConf()
	return goose.RunMigrations(goose_conf, goose_conf.MigrationsDir, version)
}
