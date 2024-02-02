package cmds

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/store"
	"golang.design/x/clipboard"
)

func GetCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get [name]",
		Short: "get password for [name] from stroe",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			err := clipboard.Init()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			service, err := NewService(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			passwd, err := service.store.Get(ctx, args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			primaryKey := GetPrimaryKey()

			decryptPasswd, err := store.Decrypt(primaryKey, []byte(passwd.CryptedPasswd))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			clipboard.Write(clipboard.FmtText, []byte(decryptPasswd))
			fmt.Printf("Name: %s; Account: %s; Url: %s; Note: %s;\n 密码已复制到剪贴板！", passwd.Name, passwd.UserName, passwd.Url, passwd.Note)
			fmt.Println(string(decryptPasswd))
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			ctx := context.Background()
			service, err := NewService(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			nameList, err := service.store.SearchName(ctx, toComplete)
			fmt.Println(nameList)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return nameList, cobra.ShellCompDirectiveKeepOrder
		},
	}

	return getCmd
}
