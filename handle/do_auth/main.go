package doauth

import (
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/wwqdrh/gokit/logger"
)

// 自定义函数

func KeyMatch(key1 string, key2 string) bool {
	i := strings.Index(key2, "*")
	if i == -1 {
		return key1 == key2
	}

	if len(key1) > i {
		return key1[:i] == key2[:i]
	}
	return key1 == key2[:i]
}

func KeyMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return (bool)(KeyMatch(name1, name2)), nil
}

var defaultReqCount = reqCount{}

type reqCount map[string]int

func (c reqCount) Add(url string) {
	c[url] += 1
}

func (c reqCount) Time(url string) int {
	return c[url]
}

func TimesMatch(reqUser, reqAPI, reqMethod string, policyUser, policyAPI, policyMethod string) bool {
	if defaultReqCount.Time(reqUser+reqAPI) >= 5 {
		return false
	}

	if reqUser != policyUser || reqAPI != policyAPI || reqMethod != policyMethod {
		return false
	}

	defaultReqCount.Add(reqUser + reqAPI)
	return true
}

func TimesMatchFunc(args ...interface{}) (interface{}, error) {
	reqUser := args[0].(string)
	reqAPI := args[1].(string)
	reqMethod := args[2].(string)
	policyUser := args[3].(string)
	policyAPI := args[4].(string)
	policyMethod := args[5].(string)

	return (bool)(TimesMatch(reqUser, reqAPI, reqMethod, policyUser, policyAPI, policyMethod)), nil
}

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
	// AddFunction 注册函数
	e.AddFunction("my_func", KeyMatchFunc)
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
	e.AddFunction("my_func", KeyMatchFunc)
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

type reqTimes struct {
	user   string
	method string
	api    string
}

func doTimes(reqs []reqTimes) []bool {
	res := make([]bool, len(reqs))

	e, err := casbin.NewEnforcer("./times/model.conf", "./times/policy.csv")
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
	}
	e.AddFunction("times", TimesMatchFunc)
	for i, item := range reqs {
		ok, err := e.Enforce(item.user, item.api, item.method)
		if err != nil {
			logger.DefaultLogger.Error(err.Error())
		}
		fmt.Printf("%v: %v\n", item, ok)
		res[i] = ok
	}
	return res
}
