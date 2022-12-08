---
marp: true
---

# Your slide deck

Start writing!

---

# 日志系统

如何搭建一个日志系统

- [Loki](https://grafana.com/docs/loki/latest/getting-started/)
- [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/)
- grafana

Loki与Promtail(loki的客户端，scra日志然后上报日志的)是grafana官方推荐的日志收集工具

通过与容器共用volume，然后promtail收集日志之后上报给loki，然后再由grafana展示处理

---

# Loki配置

> 具体查看: https://grafana.com/docs/loki/latest/configuration/examples/

```yaml
schema_config: 
```

---

# Page 3

You can directly use Windi CSS and Vue components to style and enrich your slides.

<div class="p-3">
  <Tweet id="20" />
</div>