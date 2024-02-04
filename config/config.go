package config

import (
	"fmt"
	"os"
	"path/filepath"

	c "github.com/spf13/viper"
)

func InitConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.SetConfigName("conf")                      // 配置文件的名称（不需要扩展名）
	c.AddConfigPath(filepath.Join(home, ".pm/")) // 配置文件的路径

	// 读取配置数据
	if err := c.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
