package cmds

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/config"
	"gorm.io/gorm"
)

func DelCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "del [name]",
		Short: "delete password for [name] from store",
		Args:  cobra.ExactArgs(1),
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

			_ = GetPrimaryKey()

			passwd, err := service.store.Get(cmd.Context(), args[0])
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println(err)
				os.Exit(1)
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return
			}

			var inputName string
			fmt.Printf("Are you sure you want to delete the password named %s? If yes, please enter %s: ", passwd.Name, passwd.Name)
			fmt.Scanln(&inputName)

			if inputName != passwd.Name {
				fmt.Println("input name not match")
				os.Exit(1)
			}

			if err := service.store.Delete(cmd.Context(), passwd.Name); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			config.InitConfig()
			service, err := NewService(cmd.Context())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			passwdList, err := service.store.SearchName(cmd.Context(), toComplete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			nameList := make([]string, 0)
			for _, passed := range passwdList {
				nameList = append(nameList, passed.Name)
			}
			return nameList, cobra.ShellCompDirectiveKeepOrder
		},
	}

	return getCmd
}
