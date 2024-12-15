package cmds

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/config"
	"github.com/yangchnet/pm/store"
	"gorm.io/gorm"
)

func GetCmd() *cobra.Command {
	var getCmd = &cobra.Command{
		Use:   "get [name]",
		Short: "get password for [name] from store",
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

			passwd, err := service.store.Get(cmd.Context(), args[0])
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					fmt.Println("password not found")
					os.Exit(1)
				}
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

func tabPrint(passwd *store.Passwd) {
	// Initialize a tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Write the headers
	fmt.Fprintln(w, "Name\tAccount\tUrl\tNote\tCreateTime")

	fmt.Fprintln(w, "---\t---\t---\t---\t---")

	// Write some data
	fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s", passwd.Name, passwd.UserName, passwd.Url, passwd.Note, passwd.CreateTime.Format("2006-01-02 15:04:05")))

	// Flush the writer to output the table
	w.Flush()
}
