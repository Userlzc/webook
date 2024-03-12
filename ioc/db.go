package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"project/internal/repository/dao"
)

/**
 * @Description
 * @Date 2024/3/12 19:00
 **/

func InitDB() *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/xiaohongshu?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库驱动错误")
	}
	if err := dao.InitTables(db); err != nil {
		panic("创建数据库表错误")
	}
	return db
}
