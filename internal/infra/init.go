package infra

import (
	"log"
	"os"
)

func init() {
	if err := os.MkdirAll("log", os.ModePerm); err != nil {
		log.Println("创建日志目录失败", err)
		return
	}
}
