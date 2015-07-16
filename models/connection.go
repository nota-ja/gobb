package models

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"github.com/stevenleeg/gobb/config"
	"os"
)

var db_map *gorp.DbMap

func GetDbSession() *gorp.DbMap {
	if db_map != nil {
		return db_map
	}

	db_username, _ := config.Config.GetString("database", "username")
	db_password, _ := config.Config.GetString("database", "password")
	db_database, _ := config.Config.GetString("database", "database")
	db_hostname, _ := config.Config.GetString("database", "hostname")
	db_port, _ := config.Config.GetString("database", "port")

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

	data_source := os.Getenv("DATABASE_URL")
	if data_source == "" {
		data_source = "user=" + db_username +
			" password=" + db_password +
			" dbname=" + db_database +
			" host=" + db_hostname +
			" port=" + db_port +
			" sslmode=disable"
	}

	db, err := sql.Open("postgres", data_source)

	if err != nil {
		fmt.Printf("Cannot open database! Error: %s\n", err.Error())
		return nil
	}

	db_map = &gorp.DbMap{
		Db:      db,
		Dialect: gorp.PostgresDialect{},
	}

	// TODO: Do we need this every time?
	db_map.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	db_map.AddTableWithName(Board{}, "boards").SetKeys(true, "Id")
	db_map.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")
	db_map.AddTableWithName(View{}, "views").SetKeys(false, "Id")
	db_map.AddTableWithName(Setting{}, "settings").SetKeys(true, "Key")

	return db_map
}
