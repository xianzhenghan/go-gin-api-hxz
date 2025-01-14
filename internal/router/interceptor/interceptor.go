package interceptor

import (
	"go-gin-api-hxz/internal/pkg/core"
	"go-gin-api-hxz/internal/proposal"
	"go-gin-api-hxz/internal/repository/mysql"
	"go-gin-api-hxz/internal/repository/redis"
	"go-gin-api-hxz/internal/services/admin"
	"go-gin-api-hxz/internal/services/authorized"

	"go.uber.org/zap"
)

var _ Interceptor = (*interceptor)(nil)

type Interceptor interface {
	// CheckLogin 验证是否登录
	CheckLogin(ctx core.Context) (info proposal.SessionUserInfo, err core.BusinessError)

	// CheckRBAC 验证 RBAC 权限是否合法
	CheckRBAC() core.HandlerFunc

	// CheckSignature 验证签名是否合法，对用签名算法 pkg/signature
	CheckSignature() core.HandlerFunc

	// i 为了避免被其他包实现
	i()
}

type interceptor struct {
	logger            *zap.Logger
	cache             redis.Repo
	db                mysql.Repo
	authorizedService authorized.Service
	adminService      admin.Service
}

func New(logger *zap.Logger, cache redis.Repo, db mysql.Repo) Interceptor {
	return &interceptor{
		logger:            logger,
		cache:             cache,
		db:                db,
		authorizedService: authorized.New(db, cache),
		adminService:      admin.New(db, cache),
	}
}

func (i *interceptor) i() {}
