package interceptor

import (
	"encoding/json"
	"net/http"

	"go-gin-api-hxz/configs"
	"go-gin-api-hxz/internal/code"
	"go-gin-api-hxz/internal/pkg/core"
	"go-gin-api-hxz/internal/repository/redis"
	"go-gin-api-hxz/internal/services/admin"
	"go-gin-api-hxz/pkg/errors"
	"go-gin-api-hxz/pkg/urltable"
)

func (i *interceptor) CheckRBAC() core.HandlerFunc {
	return func(c core.Context) {
		token := c.GetHeader("Token")
		if token == "" {
			c.AbortWithError(core.Error(
				http.StatusUnauthorized,
				code.AuthorizationError,
				code.Text(code.AuthorizationError)).WithError(errors.New("Header 中缺少 Token 参数")),
			)
			return
		}

		if !i.cache.Exists(configs.RedisKeyPrefixLoginUser + token) {
			c.AbortWithError(core.Error(
				http.StatusUnauthorized,
				code.CacheGetError,
				code.Text(code.CacheGetError)).WithError(errors.New("请先登录")),
			)
			return
		}

		if !i.cache.Exists(configs.RedisKeyPrefixLoginUser + token + ":action") {
			c.AbortWithError(core.Error(
				http.StatusUnauthorized,
				code.CacheGetError,
				code.Text(code.CacheGetError)).WithError(errors.New("当前账号未配置 RBAC 权限")),
			)
			return
		}

		actionData, err := i.cache.Get(configs.RedisKeyPrefixLoginUser+token+":action", redis.WithTrace(c.Trace()))
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusUnauthorized,
				code.CacheGetError,
				code.Text(code.CacheGetError)).WithError(err),
			)
			return
		}

		var actions []admin.MyActionData
		err = json.Unmarshal([]byte(actionData), &actions)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusUnauthorized,
				code.AuthorizationError,
				code.Text(code.AuthorizationError)).WithError(err),
			)
			return
		}

		if len(actions) > 0 {
			table := urltable.NewTable()
			for _, v := range actions {
				_ = table.Append(v.Method + v.Api)
			}

			if pattern, _ := table.Mapping(c.Method() + c.Path()); pattern == "" {
				c.AbortWithError(core.Error(
					http.StatusBadRequest,
					code.RBACError,
					code.Text(code.RBACError)).WithError(errors.New(c.Method() + c.Path() + " 未进行 RBAC 授权")),
				)
				return
			}
		}

	}
}
