package admin

import (
	"go-gin-api-hxz/configs"
	"go-gin-api-hxz/internal/pkg/core"
	"go-gin-api-hxz/internal/pkg/password"
	"go-gin-api-hxz/internal/repository/mysql"
	"go-gin-api-hxz/internal/repository/mysql/admin"
	"go-gin-api-hxz/internal/repository/redis"
)

func (s *service) ModifyPassword(ctx core.Context, id int32, newPassword string) (err error) {
	data := map[string]interface{}{
		"password":     password.GeneratePassword(newPassword),
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	qb := admin.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	err = qb.Updates(s.db.GetDbW().WithContext(ctx.RequestContext()), data)
	if err != nil {
		return err
	}

	s.cache.Del(configs.RedisKeyPrefixLoginUser+password.GenerateLoginToken(id), redis.WithTrace(ctx.Trace()))
	return
}
