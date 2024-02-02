package config

import (
	"log"

	c "github.com/spf13/viper"
)

func init() {
	c.SetConfigName("conf")                                 // 配置文件的名称（不需要扩展名）
	c.AddConfigPath("/home/lc/dev/github.com/yangchnet/pm") // 配置文件的路径

	err := c.ReadInConfig() // 读取配置数据
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}
