## Hi 👋

<p align="center">
  <samp>
    <a href="https://wwqdrh.github.io/blog" target="_blank">blog</a>
  </samp>
  <samp>
    <a href="https://wwqdrh.github.io/apps" target="_blank">apps</a>
  </samp>
  <samp>
    <a href="https://space.bilibili.com/538676331" target="_blank">video</a>
  </samp>
  <samp>
    <a href="https://114.132.221.80/" target="_blank">maas</a>
  </samp>
</p>

<br />

<br />

## Resource

<details>
  <summary>
    <strong>
      开发手册
    </strong>
  </summary>

- 如何做接口测试: 一个cli工具，通过json声明文件自动构造http请求并验证响应是否满足预期，提供简易的表达式语法，以及变量的上下文
  - [视频](https://www.bilibili.com/video/BV1fY411R7Dq)
  - [文档](./handle/how_api_test/README.md)
- 日志系统搭建: Loki+promtail+grafana架构,promtail与应用共享volume,避免应用在写入日志的时候直接使用网络传递日志而导致的性能开销
  - [文档](./handle/build_logcenter/README.md)
- 权限控制实践: 描述了casbin如何做权限测试，以及一个在微服务架构中，将鉴权部分移动到openresty网关中，鉴权成功才会将流量流向下游服务器
  - [文档](./handle/do_auth/README.md)
</details>

<details>
  <summary>
    <strong>
      dapr微服务平台
    </strong>
  </summary>
  
- dapr初见: dapr简单说明，以及测试官网中描述的一些简单操作
  - [视频](https://www.bilibili.com/video/BV1L24y1y75B)
  - [文档](./use_dapr/01-start/docs.md)
- 微服务功能尝试: 负载均衡、state、pubsub、secret功能测试
  - [文档](./use_dapr/02-basic/docs.md)
- docker平台下做灰度部署: 基于openresty的balancer_by_lua动态做流量的分配，实现新旧应用流量分配的动态调整
  - [视频](https://www.bilibili.com/video/BV1c84y1k79a/)
  - [文档](./use_dapr/03-gray/docs.md)
- 一个多服务系统示例: 实验dapr的sidecar，如何通信如何调用的
  - [文档](./use_dapr/04-doaapp/README.md)
- wasi与普通模式下的性能对比: 实验了一下rustserver普通模式与wasi模式下的性能对比
  - [视频](https://www.bilibili.com/video/BV1Ve4y137tW/)
  - [文档](./use_dapr/05-trywasm/README.md)
</details>
  
<br />

<br />

## Contact me

> 如果你有任何问题要问我

<p align="center">
  <img width="128" src="./know-chat.jpg" />
</p>
