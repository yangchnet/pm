package remote

import (
	"context"
	"fmt"

	"github.com/yangchnet/pm/remote/empty"
	gitremote "github.com/yangchnet/pm/remote/git"
)

type Remote interface {
	// Push 将本地的store文件推送到远端
	Push(ctx context.Context, msg ...string) error

	// Pull 将远端的store文件拉取到本地，当相关文件(夹)不存在时，应自动创建
	Pull(ctx context.Context) error

	// Init 初始化remote，返回remote配置信息
	Init(ctx context.Context) (string, error)
}

func NewRemote(ctx context.Context, remoteType string, remoteMap map[string]any) (Remote, error) {
	var remote Remote
	switch remoteType {
	case "git":
		remote = gitremote.NewGitRemote(ctx, remoteMap)
	case "empty":
		remote = empty.NewEmptyRemote()
	default:
		return nil, fmt.Errorf("未知的remote类型: %s", remoteType)
	}

	return remote, nil
}
