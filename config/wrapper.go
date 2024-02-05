package config

import (
	c "github.com/spf13/viper"
)

func Get(key string) any {
	return c.Get(key)
}

func GetString(key string) string {
	return c.GetString(key)
}

func GetStringMap(key string) map[string]any {
	return c.GetStringMap(key)
}

func GetSlice(key string) []any {
	return c.Get(key).([]any)
}

func GetStringMapString(key string) map[string]string {
	return c.GetStringMapString(key)
}
