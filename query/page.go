package query

import meta "github.com/dokidokikoi/go-common/meta/option"

type PageQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Order    string `form:"order_by"`
}

func (q *PageQuery) GetListOption() *meta.ListOption {
	return meta.NewListOption(nil, meta.WithPage(q.Page), meta.WithPageSize(q.PageSize), meta.WithOrderBy(q.Order))
}
