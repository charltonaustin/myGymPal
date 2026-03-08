package models

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
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
	return orm.RegisterDataBase("default", "postgres", dsn)
}
