package repository

import (
	"context"
	"log"
	"project/internal/domain"
	"project/internal/repository/cache"
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
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUsersRepository(repo *dao.UserDao, cache *cache.UserCache) *UsersRepository {
	return &UsersRepository{
		dao:   repo,
		cache: cache,
	}
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

func (repo *UsersRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	if err == nil {
		return du, err
	}
	// err 不为nil 就要查询数据库
	//  err有两种可能
	// 1.key 不存在 说明redis正常
	// 2.访问redis 有问题 可能是网络问题，也可能是redis本身就奔溃了
	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	du = repo.toDomain(u)
	// 异步写法
	// 另外回写缓存的时候忽略掉了错误，故需改善
	go func() {
		if err := repo.cache.Set(ctx, du); err != nil {
			// 网络崩了 或者redis崩了 缓存击穿
			log.Println(err)
		}

	}()

	return du, nil

}

func (repo *UsersRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {

	return repo.dao.UpdateById(ctx, repo.toDaoUser(user))

}

func (repo *UsersRepository) toDaoUser(user domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		Birthday: user.Birthday.UnixMilli(),
		NickName: user.NickName,
		AboutMe:  user.AboutMe,
	}

}
