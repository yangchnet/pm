package cmds

import (
	"context"
	"fmt"

	"github.com/yangchnet/pm/config"
	"github.com/yangchnet/pm/remote"
	"github.com/yangchnet/pm/remote/empty"
	gitremote "github.com/yangchnet/pm/remote/git"
	"github.com/yangchnet/pm/store"
	filestore "github.com/yangchnet/pm/store/file-store"
	sqlitestore "github.com/yangchnet/pm/store/sqlite-store"
)

type service struct {
	store  store.Store
	remote remote.Remote
}

func NewService(ctx context.Context) (*service, error) {
	remoteMap := config.GetStringMap("remote")
	if len(remoteMap) <= 0 {
		return nil, fmt.Errorf("remote not found")
	}

	remote, err := NewRemote(ctx, remoteMap["type"].(string), remoteMap)
	if err != nil {
		return nil, err
	}

	storeMap := config.GetStringMap("store")
	if len(storeMap) <= 0 {
		return nil, fmt.Errorf("store not found")
	}

	store, err := NewStore(ctx, storeMap["type"].(string), storeMap)
	if err != nil {
		return nil, err
	}

	return &service{
		store:  store,
		remote: remote,
	}, nil
}

func NewStore(ctx context.Context, storeType string, storeConfig map[string]any) (store.Store, error) {
	var localStore store.Store
	switch storeType {
	case "sqlite":
		localStore = sqlitestore.NewSqliteStore(ctx)
	case "file":
		localStore = filestore.NewFileStore(ctx)
	default:
		return nil, fmt.Errorf("未知的store类型: %s", storeType)
	}

	return localStore, nil
}

func NewRemote(ctx context.Context, remoteType string, remoteMap map[string]any) (remote.Remote, error) {
	var remote remote.Remote
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
