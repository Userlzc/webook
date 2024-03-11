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
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户不存在或密码错误")
)

type UsersService struct {
	repo *repository.UsersRepository
}

func NewUsersService(repo *repository.UsersRepository) *UsersService {
	return &UsersService{repo: repo}
}

func (svc *UsersService) SignUp(ctx context.Context, u domain.User) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(pwd)
	return svc.repo.Create(ctx, u)

}

func (svc *UsersService) Login(ctx context.Context, email string, password string) (domain.User, error) {

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

func (svc *UsersService) GetProfile(ctx context.Context, uid int64) (domain.User, error) {

	return svc.repo.FindById(ctx, uid)

}

func (svc *UsersService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)

}
