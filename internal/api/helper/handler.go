package helper

import (
	"go-gin-api-hxz/internal/pkg/core"
	"go-gin-api-hxz/internal/repository/mysql"
	"go-gin-api-hxz/internal/repository/redis"
	"go-gin-api-hxz/internal/services/authorized"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Md5 加密
	// @Tags Helper
	// @Router /helper/md5/{str} [get]
	Md5() core.HandlerFunc

	// Sign 签名
	// @Tags Helper
	// @Router /helper/sign [post]
	Sign() core.HandlerFunc
}

type handler struct {
	logger            *zap.Logger
	db                mysql.Repo
	authorizedService authorized.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:            logger,
		db:                db,
		authorizedService: authorized.New(db, cache),
	}
}

func (h *handler) i() {}
