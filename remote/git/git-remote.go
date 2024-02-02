package gitremote

import (
	"context"
	"errors"
	"fmt"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitRemote struct {
	GitUrl    string
	LocalPath string

	*GitRemoteConfig

	auth transport.AuthMethod
}

type GitRemoteConfig struct {
	Url            string `json:"url" yaml:"url"`
	Name           string `json:"name" yaml:"name"`
	Email          string `json:"email" yaml:"email"`
	PrivateKeyPath string `json:"privateKeyPath" yaml:"privateKeyPath"`
	PublicKeyPath  string `json:"publicKeyPath" yaml:"publicKeyPath"`
}

func NewGitRemote(
	ctx context.Context,
	localPath string,
	config *GitRemoteConfig,
) (*GitRemote, error) {
	auth, err := auth(config.PrivateKeyPath, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create git auth credentials: %w", err)
	}

	return &GitRemote{
		GitUrl:          config.Url,
		LocalPath:       localPath,
		GitRemoteConfig: config,
		auth:            auth,
	}, nil
}

// Push 将本地的store文件推送到远端
func (gr *GitRemote) Push(ctx context.Context) error {

	r, err := git.PlainOpen(gr.LocalPath)
	if err != nil {
		return err
	}

	if err := gr.commit(ctx, r); err != nil {
		return err
	}

	return r.Push(&git.PushOptions{
		Auth: gr.auth,
	})
}

// Pull 将远端的store文件拉取到本地
func (gr *GitRemote) Pull(ctx context.Context) error {
	r, err := git.PlainOpen(gr.LocalPath)
	if err != nil {
		return err
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	if err := w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       gr.auth,
	}); err != nil && !errors.Is(git.NoErrAlreadyUpToDate, err) {
		return err
	}

	// if errors.Is(git.NoErrAlreadyUpToDate, err) {
	// 	return nil
	// }

	return nil
}

func (gr *GitRemote) commit(ctx context.Context, g *git.Repository) error {
	w, err := g.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(".")
	if err != nil {
		return err
	}
	commitMsg := "xxx" // TODO: fixme
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gr.Name,
			Email: gr.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func auth(privateKey string, password string) (transport.AuthMethod, error) {
	return ssh.NewPublicKeysFromFile("git", privateKey, password)
}
