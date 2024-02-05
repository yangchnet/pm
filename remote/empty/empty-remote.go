package empty

import "context"

type EmptyRemote struct {
}

func NewEmptyRemote() *EmptyRemote {
	return &EmptyRemote{}
}

func (r *EmptyRemote) Init(ctx context.Context) (string, error) {
	return `remote:
  type: empty`, nil
}

func (r *EmptyRemote) Pull(ctx context.Context) error {
	return nil
}

func (r *EmptyRemote) Push(ctx context.Context, msg ...string) error {
	return nil
}
