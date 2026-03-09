package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Init() error {
	host, _ := beego.AppConfig.String("db_host")
	port, _ := beego.AppConfig.String("db_port")
	name, _ := beego.AppConfig.String("db_name")
	user, _ := beego.AppConfig.String("db_user")
	password, _ := beego.AppConfig.String("db_password")
	sslmode, _ := beego.AppConfig.String("db_sslmode")

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
