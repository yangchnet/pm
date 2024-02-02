package cmds

import (
	"context"

	"github.com/yangchnet/pm/remote"
	"github.com/yangchnet/pm/store"
)

type service struct {
	store  store.Store
	remote remote.Remote
}

func NewService(ctx context.Context) (*service, error) {
	remote, err := remote.NewRemote(ctx)
	if err != nil {
		return nil, err
	}
	return &service{
		store:  store.NewSqliteStore(ctx),
		remote: remote,
	}, nil
}
