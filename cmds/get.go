package cmds

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/config"
	"github.com/yangchnet/pm/store"
)

func GetCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get [name]",
		Short: "get password for [name] from stroe",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			config.InitConfig()

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

			if err := clipboard.WriteAll(string(decryptPasswd)); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			tabPrint(passwd)

			fmt.Println("密码已复制到剪贴板")
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			ctx := context.Background()
			config.InitConfig()
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

func tabPrint(passwd *store.Passwd) {
	// Initialize a tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Write the headers
	fmt.Fprintln(w, "Name\tAccount\tUrl\tNote")

	fmt.Fprintln(w, "---\t---\t---\t---")

	// Write some data
	fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s", passwd.Name, passwd.UserName, passwd.Url, passwd.Note))

	// Flush the writer to output the table
	w.Flush()
}
