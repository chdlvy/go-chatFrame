package db

import (
	"fmt"
	"github.com/chdlvy/go-chatFrame/pkg/common/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DBConn *gorm.DB

func NewGormDB() (*gorm.DB, error) {
	if DBConn != nil {
		return DBConn, nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Config.Mysql.Username, config.Config.Mysql.Password, config.Config.Mysql.Address, config.Config.Mysql.Database)
	fmt.Println(config.Config.Mysql.Username, config.Config.Mysql.Password, config.Config.Mysql.Address, config.Config.Mysql.Database)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	//defer sqlDB.Close()
	//设置字符集的默认值
	charset := config.Config.Mysql.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	collate := config.Config.Mysql.Collate
	if collate == "" {
		collate = "utf8mb4_unicode_ci"
	}
	//创建数据库
	sql := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s default charset %s COLLATE %s;",
		config.Config.Mysql.Database,
		charset,
		collate,
	)
	err = db.Exec(sql).Error
	if err != nil {
		return nil, fmt.Errorf("init db %w", err)
	}
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.Config.Mysql.MaxLifeTime))
	sqlDB.SetMaxOpenConns(config.Config.Mysql.MaxOpenConn)
	sqlDB.SetMaxIdleConns(config.Config.Mysql.MaxIdleConn)

	DBConn = db
	return DBConn, nil

}
