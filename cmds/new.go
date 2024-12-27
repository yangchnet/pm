package cmds

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/config"
	"github.com/yangchnet/pm/store"
	"github.com/yangchnet/pm/utils"
)

func GenerateCmd() *cobra.Command {
	var (
		account string
		note    string
		url     string
		passwd  string
		length  int32
		lower   bool // 是否使用小写字母
		upper   bool // 是否使用大写字母
		number  bool // 是否使用数字
		symbols bool // 是否使用特殊符号
	)

	var generateCmd = &cobra.Command{
		Use:   "new <name> [options]",
		Short: "generate a new password for [name]",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			config.InitConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			password := passwd
			if password == "" {
				password = utils.GeneratePassword(int(length), lower, upper, number, symbols)
			}

			primaryKey := GetPrimaryKey()

			cryptedPasswd, err := store.Encrypt(primaryKey, []byte(password))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			service, err := NewService(cmd.Context())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			passwd := store.Passwd{
				Name:          args[0],
				Url:           url,
				UserName:      account,
				Note:          note,
				CryptedPasswd: cryptedPasswd,
			}
			if err := service.store.Save(cmd.Context(), &passwd); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			clipboard.WriteAll(password)
			fmt.Println("密码已经复制到剪贴板")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			service, err := NewService(cmd.Context())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			msg := fmt.Sprintf("add password: %s", args[0])

			if err := service.remote.Push(cmd.Context(), msg); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	generateCmd.Flags().StringVar(&account, "account", "", "账户名")

	generateCmd.Flags().StringVar(&note, "note", "n", "密码的备注信息")

	generateCmd.Flags().StringVar(&url, "url", "u", "相关的url")

	generateCmd.Flags().StringVar(&passwd, "passwd", "", "已有的密码")

	generateCmd.Flags().Int32Var(&length, "length", 12, "生成的密码长度")

	generateCmd.Flags().BoolVarP(&lower, "lower", "l", false, "是否使用小写字母")
	generateCmd.Flags().BoolVarP(&upper, "upper", "u", false, "是否使用大写字母")
	generateCmd.Flags().BoolVarP(&number, "number", "n", false, "是否使用数字")
	generateCmd.Flags().BoolVarP(&symbols, "symbols", "s", false, "是否使用特殊符号")

	return generateCmd
}
