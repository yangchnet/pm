package cmds

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yangchnet/pm/config"
	"github.com/yangchnet/pm/utils"
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
	userKeyPath := config.GetString("user_key_path")
	_, err := os.Stat(userKeyPath)
	if os.IsNotExist(err) {
		return "", true
	}

	content, err := os.ReadFile(config.GetString("user_key_path"))
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

	// 如果上一次输出已经超过了给定的延时，则必须再次输出密码
	passwdLatency, err := time.ParseDuration(config.GetString("latency"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if time.Since(lastInputKeyTime) > passwdLatency {
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
	userKeyPath := config.GetString("user_key_path")

	utils.CreateDirIfNotExist(filepath.Dir(userKeyPath))

	hash := sha256.Sum256(primaryKeyByte)
	hashString := hex.EncodeToString(hash[:])

	content := []byte(hashString + " " + strconv.FormatInt(time.Now().Unix(), 10))
	if err := os.WriteFile(userKeyPath, content, 0644); err != nil {
		panic(err)
	}

	return hashString
}
