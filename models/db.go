package models

import (
	"fmt"
	"os"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func envOrConfig(envKey, configKey string) string {
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	v, _ := beego.AppConfig.String(configKey)
	return v
}

func dbConfig() (host, port, name, user, password, sslmode string) {
	host = envOrConfig("DB_HOST", "db_host")
	port = envOrConfig("DB_PORT", "db_port")
	name = envOrConfig("DB_NAME", "db_name")
	user = envOrConfig("DB_USER", "db_user")
	password = envOrConfig("DB_PASSWORD", "db_password")
	sslmode = envOrConfig("DB_SSLMODE", "db_sslmode")
	return
}

// SessionProviderDSN returns the PostgreSQL connection string for Beego's session provider.
func SessionProviderDSN() string {
	host, port, name, user, password, sslmode := dbConfig()
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		user, password, name, host, port, sslmode)
}

func Init() error {
	host, port, name, user, password, sslmode := dbConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		host, port, name, user, password, sslmode,
	)

	orm.RegisterDriver("postgres", orm.DRPostgres)
	if err := orm.RegisterDataBase("default", "postgres", dsn); err != nil {
		return err
	}

	return runMigrations(host, port, name, user, password, sslmode)
}

func runMigrations(host, port, dbname, user, password, sslmode string) error {
	pgURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)

	m, err := migrate.New("file://migrations", pgURL)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}
