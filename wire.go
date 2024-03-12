package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"project/internal/repository"
	"project/internal/repository/cache"
	"project/internal/repository/dao"
	"project/internal/service"
	"project/internal/web"
	"project/ioc"
)

/**
 * @Description
 * @Date 2024/3/12 19:02
 **/

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitSmsService,

		dao.NewUserDao,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewCacheUsersRepository,
		repository.NewCodeRepository,

		service.NewUsersService,
		service.NewCodeService,

		web.NewUserHandler,
	)
	return gin.Default()

}
