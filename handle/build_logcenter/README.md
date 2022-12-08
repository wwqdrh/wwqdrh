---
marp: true
---

# GLP(Grafana + Loki + Promtail)

---

# 日志系统

如何搭建一个日志系统

- [Loki](https://grafana.com/docs/loki/latest/getting-started/)
- [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/)
- grafana

Loki与Promtail(loki的客户端，scra日志然后上报日志的)是grafana官方推荐的日志收集工具

通过与容器共用volume，然后promtail收集日志之后上报给loki，然后再由grafana展示处理

---

# Promtail

收集日志，通过与app挂载同样的volume来共享日志目录，从而让promtail能够收集到app的日志

```yaml
# 将promtail配置成一个server，可以用于作为日志pull模式
server:
  http_listen_port: 9080
  grpc_listen_port: 0

# 配置loki的push接口
clients:
  - url: http://loki:3100/loki/api/v1/push
```

---

```yaml
# 日志读取到的位置
positions:
  filename: /tmp/positions.yaml

# 定义如何去拉取日志
scrape_configs:
- job_name: system # ui中展示的标识符
  static_configs:
  - targets:
      - localhost # 从本机拉取
    labels:
      __path__: /var/log/**/*.log # 本机上的路径
      job: varlogs # 标识符
```

---

# Loki配置

> 具体查看: https://grafana.com/docs/loki/latest/configuration/examples/

```yaml
# Enables authentication through the X-Scope-OrgID header if true
# If false, the OrgID will always be set to "fake".
auth_enabled: false

# 定义server行为，例如promtail就可以通过3100端口push日志上来
server:
  http_listen_port: 3100
  grpc_listen_port: 9096

# 全局配置，多个组件中只需要定义一处就行
common:
  path_prefix: /tmp/loki # 前端接口前缀
  storage: # 定义日志存储路径
    filesystem: # 本机路径
      chunks_directory: /tmp/loki/chunks
      rules_directory: /tmp/loki/rules
  replication_factor: 1
  ring: # 一致性hash的存储路径，用于loki cluster的
    instance_addr: 127.0.0.1
    kvstore:
      store: inmemory
```

---

```yaml
# 用于配置索引桶
schema_config:
  configs:
    - from: 2020-10-24
      # Which store to use for the index. Either aws, aws-dynamo
      # gcp, bigtable, bigtable-hashed,
      # cassandra, boltdb or boltdb-shipper.
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h
```
