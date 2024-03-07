package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

/**
 * @Description
 * @Date 2024/3/1 18:15
 **/

var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}
func (dao *UserDao) Insert(ctx context.Context, u User) error {
	// 注意数据库中唯一索引冲突的问题
	now := time.Now().Unix()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突  邮箱冲突
			return ErrDuplicateEmail
		}

	}
	return err
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err

}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
