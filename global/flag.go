package global

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
)

func NewEnforcer(url string) *casbin.Enforcer {
	adapter, err := mongodbadapter.NewAdapter(url)
	if err != nil {
		panic(err)
	}
	enforcer, err := casbin.NewEnforcer("config/policy.conf", adapter)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		panic(fmt.Sprintf("权限加载失败：%s", err.Error()))
	}
	return enforcer
}
