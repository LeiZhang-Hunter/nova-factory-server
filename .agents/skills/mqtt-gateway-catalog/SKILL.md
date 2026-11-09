---
name: mqtt-gateway-catalog
description: 在 nova-factory-server 中实现、修改或排查通过 MQTT 向边缘网关提供设备配置目录的功能，包括请求响应协议、认证、目录生成和启动配置。
---

# MQTT 网关设备配置目录

处理平台通过 MQTT 响应边缘网关设备配置目录请求的任务时，先阅读当前源码并以实际实现为准；不要将此技能中的现状当作独立于代码的协议规范。涉及跨仓库网关端时，检查对端实现并保持兼容。

## 定位代码

- `app/business/iot/daemonize/catalogwire/wire.go`：主题、请求和响应结构、ID 校验、HMAC 签名及有效期。
- `app/business/iot/daemonize/daemonizecontroller/catalog_mqtt.go`：MQTT 连接、订阅、请求处理、认证、响应和关闭清理。
- `app/business/iot/routes.go`：IoT 模块启动响应器并注册关闭函数；编译时留意 IoT build tag。
- `app/business/iot/daemonize/daemonizeservice/daemonizeserviceimpl/device_catalog.go`：从设备快照生成目录和版本。
- `app/business/iot/asset/device/devicedao/devicedaoImpl/device_catalog.go` 与 `app/business/iot/daemonize/daemonizemodels/response/`：快照查询及目录数据结构。

## 工作流程

1. 明确修改涉及协议、认证、目录数据还是运行配置，检查相关调用方、测试和工作区现有改动。修改模型或分层接口时同时遵循项目 Go 模型分层约定。
2. 保持请求主题 `nova/v1/gateways/{gatewayId}/catalog/request` 与响应主题 `nova/v1/gateways/{gatewayId}/catalog/response/{requestId}` 及 JSON 字段的兼容性；协议变更需要同步检查网关消费者。
3. 保留主题与消息中的网关身份一致性检查，以及基于网关凭据的请求认证、时间窗口和响应签名。签名基于清空 `signature` 后的 JSON 和 HMAC-SHA256；更改序列化或字段时检查双方签名兼容性。不要记录密码或密钥。
4. 目录生成需检查设备、模板、测点和连接配置的一致性。当前转换器支持 Modbus TCP 和 MQTT；按设备协议分别校验，扩展协议时同步处理数据映射和网关能力。保留大小限制、错误响应、队列边界及关闭清理等运行约束，除非任务明确要求改变。
5. 检查 `gateway_mqtt.enabled`、`broker`、`client_id`、`username`、`password` 与启动行为。未启用时响应器不连接；启用时连接或订阅失败会导致 IoT 初始化失败。需要新增配置时遵循项目现有配置方式，不写入真实凭据。

## 验证

- 根据改动运行相关测试，例如 `go test ./app/business/iot/daemonize/catalogwire`；涉及 IoT 集成时检查适用的 `iot` build tag 并运行相关包测试或构建。
- 针对行为变化验证身份、签名、时间窗口、错误响应和目录映射等实际受影响的边界。涉及 MQTT 交互时，条件允许使用隔离的测试 broker 验证订阅、请求、响应和关闭；无法运行时说明限制。
