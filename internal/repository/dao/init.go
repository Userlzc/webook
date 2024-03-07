package dao

import "gorm.io/gorm"

/**
 * @Description
 * @Date 2024/3/3 18:55
 **/

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{})

}
