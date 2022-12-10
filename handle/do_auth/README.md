# 简介

casbin权限控制实践

主要分为模型、policy两个模块

# model

## 基本语法

Model CONF 至少应包含四个部分:

- request_definition: 是访问请求的定义。 它定义了 e.Enforce(...) 函数中的参数
- policy_definition: 描述策略的含义
- policy_effect: 定义了当多个policy rule同时匹配访问请求request时,该如何对多个决策结果进行集成以实现统一决策
  - `e = some(where (p.eft == allow))`: 多条结果中只要有一个allow就可以通过
  - `e = !some(where (p.eft == deny))`: 必须全部为allow
  - `some(where (p.eft == allow)) && !some(where (p.eft == deny))`: 如果存在deny那么就deny
  - `priority(p.eft) || deny`
  - `subjectPriority(p.eft)`
- matchers
- role_definition: 用于rbac模型, 定义了角色以及继承相关的(用户和资源都可以定义)


模型CONF可以包含注释。 注释开头是 #, # 将注释该行的其余部分。

支持的模型列表

> 需要注意的是这里的模型并不是一个必须遵循的标准，而是说要完成这类功能可以如何定义，例如ACL与没有用户的ACL这些其实只是说少定义了一个变量而已，本质是没有变化的

- ACL (Access Control List, 访问控制列表)
- 具有 超级用户 的 ACL
- 没有用户的 ACL: 对于没有身份验证或用户登录的系统尤其有用。
- 没有资源的 ACL: 某些场景可能只针对资源的类型, 而不是单个资源, 诸如 write-article, read-log等权限。 它不控制对特定文章或日志的访问。
- RBAC (基于角色的访问控制)
- 支持资源角色的RBAC: 用户和资源可以同时具有角色 (或组)。
- 支持域/租户的RBAC: 用户可以为不同的域/租户设置不同的角色集。
- ABAC (基于属性的访问控制): 支持利用resource.Owner这种语法糖获取元素的属性。
- RESTful: 支持路径, 如 /res/*, /res/: id 和 HTTP 方法, 如 GET, POST, PUT, DELETE。
- 拒绝优先: 支持允许和拒绝授权, 拒绝优先于允许。
- 优先级: 策略规则按照先后次序确定优先级，类似于防火墙规则。

### RBAC

```conf
# 当然这个g也可以自定义为任意个数的参数
[role_definition]
g = _, _
g2 = _, _
```

### ABAC

支持访问结构体中的属性, 不过只能在r上使用，不能在p上使用

```conf
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == r.obj.Owner
```

## 自定义函数

在matchers中可以通过自定义函数来扩展校验

## 模式匹配

在rbac中，模式匹配能够简化g的定义

```bash
p, alice, book_group, read
g, /book/1, book_group
g, /book/2, book_group

# ==>

g, /book/:id, book_group
```

# policy

由于经常更改，所以需要使用额外的后端组件来进行存储

支持的适配器，一般就是mysql、sqlite这些存储组件

> 我一般使用gorm adapter


