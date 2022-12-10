package doauth

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/wwqdrh/gokit/logger"
)

type aclReq struct {
	sub string
	obj string
	act string
}

func doACL(reqs []aclReq) []bool {
	res := make([]bool, len(reqs))

	e, err := casbin.NewEnforcer("./acl/model.conf", "./acl/policy.csv")
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
	for i, item := range reqs {
		ok, err := e.Enforce(item.sub, item.obj, item.act)
		if err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
		fmt.Printf("%v: %v\n", item, ok)
		res[i] = ok
	}
	return res
}

type reqRBAC struct {
	sub string
	obj string
	act string
}

func doRBAC(reqs []reqRBAC) []bool {
	res := make([]bool, len(reqs))

	e, err := casbin.NewEnforcer("./rbac/model.conf", "./rbac/policy.csv")
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
	for i, item := range reqs {
		ok, err := e.Enforce(item.sub, item.obj, item.act)
		if err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
		fmt.Printf("%v: %v\n", item, ok)
		res[i] = ok
	}
	return res
}

func simple() {
	e, err := casbin.NewEnforcer("path/to/model.conf", "path/to/policy.csv")
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
	}

	sub := "alice" // 想要访问资源的用户。
	obj := "data1" // 将被访问的资源。
	act := "read"  // 用户对资源执行的操作。

	ok, err := e.Enforce(sub, obj, act)

	if err != nil {
		// 处理err
	}

	if ok == true {
		// 允许alice读取data1
	} else {
		// 拒绝请求，抛出异常
	}

	// 您可以使用BatchEnforce()来批量执行一些请求
	// 这个方法返回布尔切片，此切片的索引对应于二维数组的行索引。
	// 例如results[0] 是{"alice", "data1", "read"}的结果
	results, err := e.BatchEnforce([][]interface{}{{"alice", "data1", "read"}, {"bob", "datata2", "write"}, {"jack", "data3", "read"}})
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
	fmt.Println(results)
}
