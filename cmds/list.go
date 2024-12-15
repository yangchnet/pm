package cmds

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/config"
)

func ListCmd() *cobra.Command {
	var (
		filterString string
	)
	var listCmd = &cobra.Command{
		Use:   "list [-f <filter_string>]",
		Short: "list passwd name",
		PreRun: func(cmd *cobra.Command, args []string) {
			config.InitConfig()

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
		Run: func(cmd *cobra.Command, args []string) {
			service, err := NewService(cmd.Context())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			passwds, err := service.store.SearchName(cmd.Context(), filterString)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("%18s| %18s| %18s| %18s|\n", "name", "account", "note", "url")
			fmt.Printf("%18s| %18s| %18s| %18s|\n", "-----", "-----", "-----", "-----")
			for _, passwd := range passwds {
				fmt.Printf("%18s| %18s| %18s| %18s|\n", passwd.Name, passwd.UserName, passwd.Note, passwd.Url)
			}

		},
	}

	listCmd.Flags().StringVarP(&filterString, "filter", "f", "", "passwd filter string, default \"\"")

	return listCmd
}
