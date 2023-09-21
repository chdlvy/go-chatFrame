package orm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	if db == nil {
		dsn := "root:123456@tcp(localhost:3306)/graduation_project?charset=utf8&parseTime=True&loc=Local"
		db1, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db = db1
	}
	// 初始化表，如果存在就删了重建
	db.Migrator().DropTable(&Comment{}, &Videos{}, &Users{})
	db.AutoMigrate(&Comment{}, &Videos{}, &Users{})

	sqlDB, _ := db.DB()
	// 连接池最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 超过的连接会被连接池关闭
	sqlDB.SetMaxIdleConns(20)
}

func GetSqlConn() *gorm.DB {
	return db
}
