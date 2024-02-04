package cmds

import (
	"context"
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
		name    string
		account string
		note    string
		url     string
	)

	var generateCmd = &cobra.Command{
		Use:   "new",
		Short: "generate a new password for [name]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			config.InitConfig()
			password := generatePassword(12)

			primaryKey := GetPrimaryKey()

			cryptedPasswd, err := store.Encrypt(primaryKey, []byte(password))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			service, err := NewService(ctx)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			passwd := store.Passwd{
				Name:          name,
				Url:           url,
				UserName:      account,
				Note:          note,
				CryptedPasswd: cryptedPasswd,
			}
			if err := service.store.Save(ctx, &passwd); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			clipboard.WriteAll(password)
		},
	}

	generateCmd.Flags().StringVarP(&name, "name", "", "", "password unique name")
	generateCmd.MarkFlagRequired("name")

	generateCmd.Flags().StringVarP(&account, "account", "", "", "账户名")
	generateCmd.MarkFlagRequired("account")

	generateCmd.Flags().StringVarP(&note, "note", "", "", "备注")

	generateCmd.Flags().StringVarP(&url, "url", "", "", "相关的url")

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
