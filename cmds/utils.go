package cmds

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yangchnet/pm/config"
	"golang.org/x/term"
)

func GetPrimaryKey() string {
	primaryKey, need := needPrimaryKeyInput()
	if !need {
		return primaryKey
	}

	return inputAndSavePrimaryKey()
}

func needPrimaryKeyInput() (string, bool) {
	userKeyPath := config.GetString("userKeyPath")
	_, err := os.Stat(userKeyPath)
	if os.IsNotExist(err) {
		return "", true
	}

	content, err := os.ReadFile(config.GetString("userKeyPath"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	list := strings.Split(string(content), " ")
	if len(list) != 2 {
		fmt.Println("user key file format error")
		os.Exit(1)
	}

	timestamp, err := strconv.ParseInt(list[1], 10, 64)
	if err != nil {
		fmt.Println("user key file format error")
		os.Exit(1)
	}

	lastInputKeyTime := time.Unix(timestamp, 0)

	// 如果上一次输出已经超过了24小时，则必须再次输出密码
	if time.Now().Sub(lastInputKeyTime) > time.Hour*24 {
		_ = os.Remove(userKeyPath)
		return "", true
	}

	return list[0], false
}

func inputAndSavePrimaryKey() string {
	fmt.Println("Please input your primary password:")
	primaryKeyByte, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return savePrimaryKey(primaryKeyByte)
}

func savePrimaryKey(primaryKeyByte []byte) string {
	filePath := config.GetString("userKeyPath")

	// 检查文件路径是否存在
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 如果文件路径不存在，创建它
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(err)
		}
	}

	hash := sha256.Sum256(primaryKeyByte)
	hashString := hex.EncodeToString(hash[:])

	content := []byte(hashString + " " + strconv.FormatInt(time.Now().Unix(), 10))
	if err := os.WriteFile(config.GetString("userKeyPath"), content, 0644); err != nil {
		panic(err)
	}

	return hashString
}

func extractGitInfo(gitURL string) (string, string, string, error) {
	re := regexp.MustCompile(`^git@(.*):(.*)/(.*).git$`)
	matches := re.FindStringSubmatch(gitURL)

	if len(matches) != 4 {
		return "", "", "", fmt.Errorf("git url格式不对")
	}

	return matches[1], matches[2], matches[3], nil
}
