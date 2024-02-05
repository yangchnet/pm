package cmds

import (
	"context"
	"fmt"

	"github.com/yangchnet/pm/config"
	"github.com/yangchnet/pm/remote"
	"github.com/yangchnet/pm/store"
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

	remote, err := remote.NewRemote(ctx, remoteMap["type"].(string), remoteMap)
	if err != nil {
		return nil, err
	}
	return &service{
		store:  store.NewSqliteStore(ctx),
		remote: remote,
	}, nil
}
