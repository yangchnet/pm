package cmds

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/config"
)

func PushCmd() *cobra.Command {

	var pushCmd = &cobra.Command{
		Use:   "push",
		Short: "push passwd store to remote",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			config.InitConfig()

			service, err := NewService(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if err := service.remote.Push(ctx); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	return pushCmd

}
