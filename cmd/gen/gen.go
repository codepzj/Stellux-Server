package main

import (
	"flag"

	"github.com/codepzj/Stellux-Server/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "../../conf/dev.yaml", "config path, eg: -conf dev.yaml")
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../../internal/domain/query",                 // 生成代码的输出路径
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface, // 生成模式
		FieldNullable: true,                                          // 数据库字段设置为null值的列，在更新或插入的时候设置为nil，而不是零值
	})

	cfg := conf.GetConfig(flagconf)

	db, _ := gorm.Open(mysql.Open(cfg.Mysql.Dsn), &gorm.Config{})
	g.UseDB(db)

	g.ApplyBasic(
		g.GenerateAllTable()..., // 生成所有表
	)
	// Generate the code
	g.Execute()
}
