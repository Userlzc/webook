package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"project/internal/domain"
	"project/internal/repository"
)

/**
 * @Description
 * @Date 2024/3/1 17:49
 **/

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUser
	ErrInvalidUserOrPassword = errors.New("用户不存在或密码错误")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	GetProfile(ctx context.Context, uid int64) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}
type usersService struct {
	repo repository.UserRepository
}

func NewUsersService(repo repository.UserRepository) UserService {
	return &usersService{repo: repo}
}

func (svc *usersService) SignUp(ctx context.Context, u domain.User) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(pwd)
	return svc.repo.Create(ctx, u)

}

func (svc *usersService) Login(ctx context.Context, email string, password string) (domain.User, error) {

	u, err := svc.repo.FindByEmail(ctx, email)

	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err

	}
	// 对密码的一致性进行判断
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil

}

func (svc *usersService) GetProfile(ctx context.Context, uid int64) (domain.User, error) {

	return svc.repo.FindById(ctx, uid)

}

func (svc *usersService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)

}

func (svc *usersService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// 那就说明用户没有找到，需要创建
	err = svc.repo.Create(ctx, u)
	// 两种可能 一种是唯一索引冲突 另一种恰好是系统错误

	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	// 要么err==nil 要么ErrDuplicateUser 也代表用户存在
	// 主从延迟  理论上讲强制走主库
	return svc.repo.FindByPhone(ctx, phone)

}
