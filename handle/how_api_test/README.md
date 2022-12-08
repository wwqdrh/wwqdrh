> 实验地址: https://github.com/wwqdrh/brodcast/tree/main/how_api_test

## 接口测试

接口的正确性，比起使用postman或者全丢给测试同学，自己随着代码版本一起维护测试用例更稳妥。

不过如果你是多语言者，不同的语言维护着不同的测试代码，其实也是一种消耗精力的事。

所以写了个小工具，支持验证接口行为是否满足json描述文件的期望

- 变量上下文
- 简易的表达式
- 文件上传

## 语法列表

- `$res.$body.$json`
- `$env.a = 1`
- `$env.token = $res.$body.$json.accessToken`
- `@contain($res.$body.$str, "ok")`
- `$env.token`
- `$file:./testdata/avatar.png`
- `{{ accessToken }}`

## 示例

```json
[
    {
        "name": "用户注册",
        "url": "http://127.0.0.1:8000/api/user/register",
        "method": "post",
        "data": {
            "name": "wwqdrh",
            "gender": "1",
            "mobile": "123456789",
            "password": "123456"
        },
        "content-type": "application/json",
        "expect": [
            "@contain($res.$body.$str, \"ok\")"
        ]
    },
    {
        "name": "用户登录",
        "url": "http://127.0.0.1:8000/api/user/login",
        "method": "post",
        "data": {
            "name": "wwqdrh",
            "gender": "1",
            "mobile": "123456789",
            "password": "123456"
        },
        "content-type": "application/json",
        "expect": [
            "@contain($res.$body.$str, \"ok\")"
        ],
        "event": [
            "$env.token = $res.$body.$json.accessToken"
        ]
    },
    {
        "name": "用户信息",
        "url": "http://127.0.0.1:8000/api/user/userinfo",
        "method": "post",
        "data": {
            "name": "wwqdrh",
            "gender": "1",
            "mobile": "123456789",
            "password": "123456"
        },
        "content-type": "application/json",
        "header": [
            "Authorization: bearer {{ token }}"
        ],
        "expect": [
            "@contain($res.$body.$str, \"ok\")"
        ]
    }
]
```
