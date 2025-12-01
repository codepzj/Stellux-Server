package conf

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	MongoDB MongoDB `mapstructure:"MongoDB"`
	Server  Server  `mapstructure:"Server"`
}

type MongoDB struct {
	Host                    string `mapstructure:"HOST"`
	Port                    int    `mapstructure:"PORT"`
	MongoInitdbRootUsername string `mapstructure:"MONGO_INITDB_ROOT_USERNAME"`
	MongoInitdbRootPassword string `mapstructure:"MONGO_INITDB_ROOT_PASSWORD"`
	MongoInitdbDatabase     string `mapstructure:"MONGO_INITDB_DATABASE"`
	MongoUsername           string `mapstructure:"MONGO_USERNAME"`
	MongoPassword           string `mapstructure:"MONGO_PASSWORD"`
}

type Server struct {
	Port      int    `mapstructure:"PORT"`
	JwtSecret string `mapstructure:"JWT_SECRET"`
	Log       Log    `mapstructure:"Log"`
}

type Log struct {
	Env        string `mapstructure:"ENV"`
	Level      string `mapstructure:"LEVEL"`
	File       string `mapstructure:"FILE"`
	ErrorFile  string `mapstructure:"ERROR_FILE"`
	MaxSize    int    `mapstructure:"MAX_SIZE"`
	MaxBackups int    `mapstructure:"MAX_BACKUPS"`
	MaxAge     int    `mapstructure:"MAX_AGE"`
	Compress   bool   `mapstructure:"COMPRESS"`
}

func GetConfig(cfgPath string) *Config {
	v := viper.New()

	// 设置配置文件路径
	v.SetConfigFile(cfgPath)

	// 读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		log.Println("viper读取配置文件失败", err)
		panic("viper读取配置文件失败")
	}

	//监控配置文件变化
	v.WatchConfig()

	v.OnConfigChange(func(in fsnotify.Event) {
		if err := v.Unmarshal(&Config{}); err != nil {
			log.Println("viper读取配置文件失败", err)
		}
	})

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Println("viper读取配置文件失败", err)
	}

	// 输出配置
	log.Println("配置信息", "config", config)

	return &config
}
