package database

import (
	"fmt"
	"sync"

	"gitee.com/openeuler/ci-bot/pkg/cibot/config"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"

	// import mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DBConnection is used to operate on database
var DBConnection *gorm.DB

// once is used to init DBConnection
var once sync.Once

// New Database Connection
func New(config config.Config) error {
	once.Do(func() {
		// Connect DataBase
		dbConnection, err := ConnectDataBase(config)
		if err != nil {
			glog.Errorf("connecting to database error: %v", err)
			panic(err)
		}
		err = UpgradeDataBase(dbConnection)
		if err != nil {
			glog.Errorf("upgrading database error: %v", err)
			panic(err)
		}
		DBConnection = dbConnection
	})
	return nil
}

// ConnectDataBase connect database
func ConnectDataBase(config config.Config) (*gorm.DB, error) {
	// connect to database
	connStr := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
		config.DataBaseUserName,
		config.DataBasePassword,
		config.DataBaseHost,
		config.DataBasePort,
		config.DataBaseName)
	glog.Infof("connecting str: %v", connStr)

	return gorm.Open(config.DataBaseType, connStr)
}

// UpgradeDataBase upgrades tables and datas
func UpgradeDataBase(db *gorm.DB) error {

	// upgrades defines
	upgrades := make([]func() error, 1)
	upgrades[0] = func() error {
		// table upgrades
		if err := db.Exec(UpgradesTableSQL).Error; err != nil {
			return err
		}
		// table cla_details
		if err := db.Exec(CLADetailsTableSQL).Error; err != nil {
			return err
		}
		// table projectfiles
		if err := db.Exec(ProjectFilesTableSQL).Error; err != nil {
			return err
		}
		// table repositories
		if err := db.Exec(RepositoriesTableSQL).Error; err != nil {
			return err
		}
		// table privileges
		if err := db.Exec(PrivilegesTableSQL).Error; err != nil {
			return err
		}

		return nil
	}

	// Get UpgradeID
	var lastUpgrade = -1
	if db.HasTable(UpgradesTableName) {
		var ups []Upgrades
		err := db.Order("upgrade_id desc").Find(&ups).Error
		if err != nil {
			glog.Errorf("getting upgrades error: %v", err)
			return err
		}

		// Get value
		if len(ups) > 0 {
			lastUpgrade = ups[0].UpgradeID
		}
	}

	// Exec upgrades one by one
	for index := lastUpgrade + 1; index < len(upgrades); index++ {
		// Begin transaction
		tx := db.Begin()

		// Exec upgrades
		err := upgrades[index]()
		if err != nil {
			tx.Rollback()
			return err
		}

		// Save last upgrade
		err = db.Save(&Upgrades{
			UpgradeID: index,
		}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// End transaction
		tx.Commit()
	}

	return nil
}
