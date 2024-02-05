package cmds

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/yangchnet/pm/config"
	"github.com/yangchnet/pm/store"
)

func GenerateCmd() *cobra.Command {
	var (
		account string
		note    string
		url     string
	)

	var generateCmd = &cobra.Command{
		Use:   "new",
		Short: "generate a new password for [name]",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			config.InitConfig()
		},
		Run: func(cmd *cobra.Command, args []string) {
			password := generatePassword(12)

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

	generateCmd.Flags().StringVarP(&account, "account", "a", "", "账户名")

	generateCmd.Flags().StringVarP(&note, "note", "n", "", "密码的备注信息")

	generateCmd.Flags().StringVarP(&url, "url", "u", "", "相关的url")

	return generateCmd
}

const (
	// 小写英文字母
	characters = "abcdefghijklmnopqrstuvwxyz"

	// 大写英文字母
	upperCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// 数字
	nums = "0123456789"

	// 特殊符号
	symbols = "~!@#$%^&*()_+`-={}|[]:<>?,./"
)

func generatePassword(length int) string {
	charsets := shuffleString(characters + upperCharacters + nums + symbols)

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charsets[seededRand.Intn(len(charsets))]
	}
	return string(b)
}

func shuffleString(s string) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	runes := []rune(s)
	r.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}
