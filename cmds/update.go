package cmds

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/config"
	"github.com/yangchnet/pm/store"
	"gorm.io/gorm"
)

func UpdateCmd() *cobra.Command {
	var (
		username string
		password string
		url      string
		note     string
	)
	var updateCmd = &cobra.Command{
		Use:   "update <password> [-u username -p password -l url -n note]",
		Short: "update passwd",
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

			if len(args) < 1 {
				fmt.Println("必须执行密码名称!")
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
					fmt.Printf("不存在的密码 [%s]", args[0])
					os.Exit(1)
				}
				os.Exit(0)
			}

			if password != "" {
				primaryKey := GetPrimaryKey()
				passwd.CryptedPasswd, err = store.Encrypt(primaryKey, []byte(password))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			if note != "" {
				passwd.Note = note
			}

			if url != "" {
				passwd.Url = url
			}

			if username != "" {
				passwd.UserName = username
			}

			fmt.Printf("note: %s; url: %s", note, url)

			err = service.store.Update(cmd.Context(), args[0], passwd)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	updateCmd.Flags().StringVarP(&username, "username", "u", "", "passwd username")
	updateCmd.Flags().StringVarP(&url, "url", "l", "", "passwd url")
	updateCmd.Flags().StringVarP(&note, "note", "n", "", "passwd note")
	updateCmd.Flags().StringVarP(&password, "password", "p", "", "raw passwd")

	return updateCmd
}
