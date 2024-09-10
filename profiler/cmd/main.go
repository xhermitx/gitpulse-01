package cmd

import (
	"log"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/xhermitx/gitpulse-01/profiler/cmd/server"
	"github.com/xhermitx/gitpulse-01/profiler/configs"
	"github.com/xhermitx/gitpulse-01/profiler/service/cache"
	"github.com/xhermitx/gitpulse-01/profiler/service/candidate"
	"github.com/xhermitx/gitpulse-01/profiler/service/git"
	"github.com/xhermitx/gitpulse-01/profiler/service/queue"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := newMySQLStorage(mysql.Config{
		DSNConfig: &mysqlCfg.Config{
			User:                 configs.Envs.DBUser,
			Passwd:               configs.Envs.DBPassword,
			DBName:               configs.Envs.DBName,
			Addr:                 configs.Envs.DBAddress,
			Net:                  "tcp",
			AllowNativePasswords: true,
			ParseTime:            false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	ch, err := queue.RMQConnect(5, configs.Envs.RabbitMQAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Cache
	rdb := redis.NewClient(&redis.Options{
		Addr:     configs.Envs.RedisAddr,
		Password: "",
		DB:       0,
	})

	s := candidate.NewCandidateStore(db)
	g := git.NewGitService()
	r := queue.NewRabbitMQClient(ch)
	c := cache.NewRedisClient(rdb)

	svr := server.NewServer(s, g, r, c)
	svr.Run()
}

func newMySQLStorage(cfg mysql.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(cfg), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}

	return db, nil
}
