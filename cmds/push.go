package cmds

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func PushCmd() *cobra.Command {

	var pushCmd = &cobra.Command{
		Use:   "push",
		Short: "push passwd store to remote",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

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
