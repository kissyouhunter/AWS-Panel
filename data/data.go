package data

import (
	"fmt"
	"github.com/Yuzuki616/Aws-Panel/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func DbInit(path string) error {
	if utils.IsNotFound(path) {
		defer func() {
			log.Info("Create default admin")
			err := CreateUser("admin", "admin@admin.com", "admin", 1)
			if err != nil {
				log.Error("Create user error: ", err)
				return
			}
			log.Info("Done. account: admin password: admin123456")
		}()
	}
	db, openErr := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if openErr != nil {
		return fmt.Errorf("open db error: %v", openErr)
	}
	err := db.AutoMigrate(UserData{}, AwsSecret{})
	if err != nil {
		return fmt.Errorf("AutoMigrate error: %v", err)
	}
	Db = db
	return nil
}
