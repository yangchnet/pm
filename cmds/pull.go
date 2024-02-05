package cmds

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/config"
)

func PullCmd() *cobra.Command {
	var pushCmd = &cobra.Command{
		Use:   "pull",
		Short: "pull passwd store from remote",
		PreRun: func(cmd *cobra.Command, args []string) {
			config.InitConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {

			service, err := NewService(cmd.Context())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if err := service.remote.Pull(cmd.Context()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	return pushCmd
}
