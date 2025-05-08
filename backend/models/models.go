package models

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"learning.com/learning/utils/helper"
)

type Models struct {
	Users  UserModel
	Tokens TokenModel
}

func NewModels(cfg *DatabaseConfig, helper *helper.Helper) (*Models, error) {
	db, err := sqlx.Open("postgres", cfg.Dsn())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.MaxIdleDuration())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx) // 测试数据库连接
	if err != nil {
		return nil, err
	}

	return &Models{
		Users:  UserModel{DB: db, helper: helper},
		Tokens: TokenModel{DB: db},
	}, nil
}
