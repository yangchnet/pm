package cmds

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/config"
)

func checkIfRepositioryExists() bool {
	storePath := config.GetString("local.path")
	_, _, repo, err := extractGitInfo(config.GetString("remote.git.url"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = git.PlainOpen(filepath.Join(storePath, repo))
	if errors.Is(err, git.ErrRepositoryNotExists) {
		return false
	}

	return true
}

func PullCmd() *cobra.Command {
	var pushCmd = &cobra.Command{
		Use:   "pull",
		Short: "pull passwd store from remote",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			config.InitConfig()

			service, err := NewService(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if !checkIfRepositioryExists() {
				if err := service.remote.Init(ctx); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			if err := service.remote.Pull(ctx); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	return pushCmd
}
