package main

import (
	"fmt"
	"log"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/xhermitx/gitpulse-01/backend/cmd/api"
	"github.com/xhermitx/gitpulse-01/backend/config"
	"github.com/xhermitx/gitpulse-01/backend/db"
	"gorm.io/driver/mysql"
)

func main() {

	db, err := db.NewMySQLStorage(mysql.Config{
		DSNConfig: &mysqlCfg.Config{
			User:                 config.Envs.DBUser,
			Passwd:               config.Envs.DBPassword,
			DBName:               config.Envs.DBName,
			Addr:                 config.Envs.DBAddress,
			Net:                  "tcp",
			AllowNativePasswords: true,
			ParseTime:            true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(fmt.Sprintf("%s:%s", config.Envs.PublicHost, config.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
