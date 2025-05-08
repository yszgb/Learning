package backgroud

import "learning.com/learning/utils/jsonlog"

// 后台服务

type Background struct {
}

func New(logger *jsonlog.Logger) *Background {
	return &Background{}
}
