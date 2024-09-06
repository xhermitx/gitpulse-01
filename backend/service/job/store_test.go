package job

import (
	"log"
	"testing"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/xhermitx/gitpulse-01/backend/config"
	"github.com/xhermitx/gitpulse-01/backend/db"
	"gorm.io/driver/mysql"
)

const jobId = "7da6a7cc-08cc-4aeb-8895-63185ccfcf5e"

func TestDeleteJob(t *testing.T) {

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

	jobStore := Store{
		db: db,
	}

	// t.Run("Expected to delete a job without error", func(t *testing.T) {
	// 	err := jobStore.DeleteJob(jobId)
	// 	if err != nil {
	// 		t.Errorf("Failed to delete Job")
	// 	}
	// })

	t.Run("Expected to return a job", func(t *testing.T) {
		job, err := jobStore.FindJobById(jobId)
		if err != nil || job == nil {
			t.Error("Failed to fetch the job")
		}
	})
}
