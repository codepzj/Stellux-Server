package global

import (
	"log/slog"

	"github.com/spf13/viper"
)

func InitViper(config string) {
	viper.SetConfigFile(config)
	slog.Info("使用配置文件", "config", config)
	err := viper.ReadInConfig()
	if err != nil {
		panic("viper读取配置文件失败")
	}
}
