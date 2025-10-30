package conf

import (
	"log/slog"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Mysql  *Mysql  `mapstructure:"mysql"`
	Log    *Log    `mapstructure:"log"`
	Server *Server `mapstructure:"server"`
}

// MySQL配置
type Mysql struct {
	Dsn              string `mapstructure:"dsn"`                // 数据库连接字符串
	LogType          int    `mapstructure:"log_type"`           // 0 控制台 1 文件
	LogFile          string `mapstructure:"log_file"`           // 日志文件路径
	SlowSqlThreshold int    `mapstructure:"slow_sql_threshold"` // 慢查询阈值，单位秒
	LogLevel         string `mapstructure:"log_level"`          // 日志级别 silent info warn error
}

// Log日志配置
type Log struct {
	Level            string `mapstructure:"level"`              // 日志级别
	FullLogFilename  string `mapstructure:"full_log_filename"`  // 日志文件名
	ErrorLogFilename string `mapstructure:"error_log_filename"` // 错误日志文件名
	MaxSize          int    `mapstructure:"max_size"`           // 单个日志文件最大尺寸，单位MB
	MaxBackups       int    `mapstructure:"max_backups"`        // 最大保留日志文件数量
	MaxAge           int    `mapstructure:"max_age"`            // 最大保留天数
	Compress         bool   `mapstructure:"compress"`           // 是否压缩日志文件
}

// Server配置
type Server struct {
	Port      int    `mapstructure:"port"`
	JwtSecret string `mapstructure:"jwt_secret"`
}

// 获取配置
func GetConfig(cfgPath string) *Config {
	v := viper.New()

	// 设置配置文件路径
	v.SetConfigFile(cfgPath)

	// 读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		slog.Error("viper读取配置文件失败", "error", err.Error())
		panic("viper读取配置文件失败")
	}

	//监控配置文件变化
	v.WatchConfig()

	v.OnConfigChange(func(in fsnotify.Event) {
		if err := v.Unmarshal(&Config{}); err != nil {
			slog.Error(err.Error())
		}
	})

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		slog.Error(err.Error())
	}

	// 输出配置
	slog.Info("配置信息", "config", config)

	return &config
}
