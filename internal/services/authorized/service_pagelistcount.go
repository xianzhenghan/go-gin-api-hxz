package authorized

import (
	"go-gin-api-hxz/internal/pkg/core"
	"go-gin-api-hxz/internal/repository/mysql"
	"go-gin-api-hxz/internal/repository/mysql/authorized"
)

func (s *service) PageListCount(ctx core.Context, searchData *SearchData) (total int64, err error) {
	qb := authorized.NewQueryBuilder()
	qb = qb.WhereIsDeleted(mysql.EqualPredicate, -1)

	if searchData.BusinessKey != "" {
		qb.WhereBusinessKey(mysql.EqualPredicate, searchData.BusinessKey)
	}

	if searchData.BusinessSecret != "" {
		qb.WhereBusinessSecret(mysql.EqualPredicate, searchData.BusinessSecret)
	}

	if searchData.BusinessDeveloper != "" {
		qb.WhereBusinessDeveloper(mysql.EqualPredicate, searchData.BusinessDeveloper)
	}

	total, err = qb.Count(s.db.GetDbR().WithContext(ctx.RequestContext()))
	if err != nil {
		return 0, err
	}

	return
}
