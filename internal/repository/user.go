package repository

import (
	"context"
	"project/internal/domain"
	"project/internal/repository/dao"
)

/**
 * @Description
 * @Date 2024/3/1 17:49
 **/

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UsersRepository struct {
	dao *dao.UserDao
}

func NewUsersRepository(repo *dao.UserDao) *UsersRepository {
	return &UsersRepository{dao: repo}
}

func (repo *UsersRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

}

func (repo *UsersRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil

}
func (repo *UsersRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}
