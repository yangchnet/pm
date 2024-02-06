package store

import (
	"context"
)

type Store interface {
	// Init 初始化存储
	Init(ctx context.Context) (string, error)

	// Save 在使用cryptFunc对密码密文进行存储
	Save(ctx context.Context, passwd *Passwd) error

	// Get 获取密码
	Get(ctx context.Context, name string) (*Passwd, error)

	// SearchName 根据名称进行搜索并给出名称列表
	SearchName(ctx context.Context, name string) ([]string, error)

	// Delete 删除一个记录
	Delete(ctx context.Context, name string) error
}
