# confighotload 热加载客户端

通过 gRPC 双向流实时接收配置中心下发的 agent 配置变更，并在不重启服务的情况下将变更应用到本地数据库。

## 目录结构

```
confighotload/
├── runner/                  # 客户端运行时实现
│   ├── run.go              # 主入口：New / Run / runSession / heartbeatLoop / recvLoop
│   ├── types.go            # 核心接口和类型定义
│   ├── processor.go        # ProcessorRegistry（策略模式的管理器）
│   ├── agent_processor.go  # 默认 agent 配置处理器（解码 + 写库）
│   └── storage.go          # 存储适配器（对接现有 store.Storage）
├── v1/                     # protobuf 生成的 gRPC 代码（只读）
│   ├── agent.pb.go
│   └── agent_grpc.pb.go
└── agent.proto            # gRPC 服务定义（来源文件）
```

## 核心架构

```
┌─────────────────────────────────────────────────────────────────┐
│                          Runner                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  runSession: 建立连接 → bootstrap → 双向流 + 心跳并发的会话 │  │
│  └───────────────────────────────────────────────────────────┘  │
│                            │                                    │
│                            ▼                                    │
│  ┌──────────────────── ProcessorRegistry ────────────────────┐  │
│  │           Match(change) → 第一个 CanProcess=true 的 Processor│  │
│  └──────────────────────────┬────────────────────────────────┘  │
│                             ▼                                     │
│              ┌──────── AgentConfigProcessor ────────┐            │
│              │  1. Decoder.Decode(content)           │            │
│              │  2. AgentRepository.ApplyHotloadPatch │            │
│              └──────────────┬──────────────────────────┘        │
│                             ▼                                     │
│              ┌──── StorageAgentRepository ──────────┐           │
│              │  storage.GetAgentDao().ApplyHotloadPatch │        │
│              └──────────────┬──────────────────────────┘        │
│                             ▼                                     │
│              ┌─────────── AgentDao.ApplyHotloadPatch ─────────┐│
│              │          UPDATE ai_agents SET ... WHERE ...     ││
│              └──────────────┬──────────────────────────────────┘│
│                             ▼                                     │
│                    ┌──────── mysql ─────────┐                   │
│                    │     ai_agents 表        │                   │
│                    └────────────────────────┘                   │
└─────────────────────────────────────────────────────────────────┘
```

## gRPC 服务定义

参见 `agent.proto`：

```protobuf
service AgentControllerService {
  // agent 心跳，上报当前持有的 ConfigUuid，服务端返回最新版本
  rpc AgentHeartbeat(AgentHeartbeatReq) returns (AgentHeartbeatRes) {}

  // 获取指定版本的完整配置内容
  rpc AgentGetConfig(AgentGetConfigReq) returns (AgentGetConfigRes) {}

  // 双向流订阅配置变更
  rpc WatchAgentChanges(stream WatchAgentRequest) returns (stream ConfigChangeEvent);
}
```

- **AgentHeartbeat**：每 30 秒调用一次，上报 `AgentId` + `ConfigUuid`；若服务端版本不一致，立即触发全量同步
- **AgentGetConfig**：按版本 UUID 拉取完整的 JSON 配置内容
- **WatchAgentChanges**：客户端发送订阅请求（带 `GatewayId`），之后持续接收 `ConfigChangeEvent` 推送

## 设计模式

| 模式 | 应用位置 | 作用 |
|------|---------|------|
| **策略模式** | `Processor` 接口 + `ProcessorRegistry` | 不同类型的配置（agent / file / builtin）由不同 Processor 处理，按需扩展 |
| **适配器模式** | `AgentRepository` 接口 + `StorageAgentRepository` | 将 Runner 的接口解耦为不依赖具体 ORM 的抽象，通过 store.Storage 复用现有 DAO |
| **观察者模式** | `Observer` 接口 | 配置应用成功后回调，用于缓存失效、指标上报、下游 reload |
| **装饰器模式（准备）** | `grpc.DialOptions` | 可叠加重试、熔断、日志、Tracing 等装饰器 |

## 核心接口

### Runner 入口

```go
// 方式一：一步启动（内部自动构建 agent 配置处理器）
func Run(ctx context.Context, cfg Config) error

// 方式二：精细控制
r, err := New(cfg)
r.Run(ctx)
state := r.State() // 用于健康检查
```

### Config 参数说明

```go
type Config struct {
    Endpoint          string            // gRPC 服务地址，必填
    GatewayID         int64             // 网关 ID，必填
    AgentID           int64             // agent 实例 ID，必填
    InitialConfigUUID string            // 初始版本 UUID，可为空
    HeartbeatInterval time.Duration      // 心跳间隔，默认 30s
    SyncTimeout       time.Duration      // gRPC 调用超时，默认 10s
    Backoff           RetryBackoff      // 重连退避，默认 {1s, 30s}
    DialOptions       []grpc.DialOption // gRPC 连接选项，默认 insecure
    Logger            *log.Logger       // 日志，默认 confighotload 子日志
    Registry          *ProcessorRegistry // 处理器注册中心，默认自动构建
    AgentRepository   AgentRepository   // 持久化适配器（Registry 为空时必填）
    Decoder           AgentPatchDecoder // JSON 解码器，默认 JSONAgentPatchDecoder
    Observers         []Observer        // 配置已应用回调，可选
}
```

### Processor 接口（策略模式）

```go
type Processor interface {
    Name() string
    CanProcess(change ChangeEnvelope) bool // 匹配规则
    Process(ctx context.Context, change ChangeEnvelope) (ApplyResult, error)
}
```

### Observer 接口（观察者模式）

```go
type Observer interface {
    // 变更应用成功后回调
    // result 中包含实际处理的 Processor 名称和应用版本
    OnConfigApplied(ctx context.Context, change ChangeEnvelope, result ApplyResult) error
}
```

### AgentPatchDecoder 接口（解码适配）

```go
type AgentPatchDecoder interface {
    Decode(content string) (models.AgentHotloadPatch, error)
}
```

默认实现 `JSONAgentPatchDecoder` 特性：
- 兼容 camelCase 和 snake_case 字段命名
- 对于 `mcpServerIds` 等数组类型字段，自动序列化为 JSON 字符串以适配数据库 TEXT 存储

## 启动示例

```go
import (
    "context"
    "nova-ai-agent/internal/services/confighotload/runner"
)

func startHotload(ctx context.Context) error {
    cfg := runner.Config{
        Endpoint:    "dns:///(agent-controller:8080)",
        GatewayID:   1,
        AgentID:     100,
        Storage:     storage, // *store.Storage
    }

    // 添加配置应用后的回调（如缓存失效）
    cfg.Observers = append(cfg.Observers, &MyCacheInvalidator{})

    return runner.Run(ctx, cfg)
}
```

## 扩展方式

### 1. 自定义 Processor（处理其他配置类型）

```go
type FileConfigProcessor struct{}

func (p *FileConfigProcessor) Name() string { return "file_config" }

func (p *FileConfigProcessor) CanProcess(change runner.ChangeEnvelope) bool {
    return change.RefType == "file"
}

func (p *FileConfigProcessor) Process(ctx context.Context, change runner.ChangeEnvelope) (runner.ApplyResult, error) {
    // 写入本地配置文件
    return runner.ApplyResult{Processor: p.Name(), Version: change.Version}, nil
}

// 注册
registry := runner.NewProcessorRegistry(
    runner.NewAgentConfigProcessor(agentID, repo, decoder),
)
registry.Register(&FileConfigProcessor{})
```

### 2. 自定义 Decoder（适配其他 JSON 格式）

```go
type MyDecoder struct{}

func (d MyDecoder) Decode(content string) (models.AgentHotloadPatch, error) {
    // 自定义解析逻辑
    return models.AgentHotloadPatch{}, nil
}

cfg.Decoder = MyDecoder{}
```

### 3. 自定义 Observer（缓存失效、指标上报）

```go
type CacheInvalidator struct{}

func (o *CacheInvalidator) OnConfigApplied(ctx context.Context, change runner.ChangeEnvelope, result runner.ApplyResult) error {
    // 清除 viper 缓存
    viper.Reset()
    return nil
}

cfg.Observers = append(cfg.Observers, &CacheInvalidator{})
```

### 4. 自定义 AgentRepository（切换存储后端）

```go
type RedisAgentRepository struct{}

func (r *RedisAgentRepository) ApplyHotloadPatch(ctx context.Context, agentID int64, patch models.AgentHotloadPatch) error {
    // 写入 Redis
    return nil
}

cfg.AgentRepository = &RedisAgentRepository{}
```

## 配置变更来源追踪

变更通过 `ChangeSource` 标识触发来源：

| 来源 | 触发时机 | 说明 |
|------|---------|------|
| `bootstrap` | `runSession` 初始化时首次 heartbeat 版本对比 | 启动时立即同步 |
| `heartbeat` | 心跳协程发现 ConfigUuid 不一致 | 运行时周期性检查 |
| `stream` | `WatchAgentChanges` 收到推送 | 实时推送 |

可用于日志分级、指标标签区分。

## 状态快照

通过 `Runner.State()` 获取运行时状态，用于健康检查：

```go
type StateSnapshot struct {
    Connected      bool      // 是否已连接
    ConfigUUID     string    // 当前持有的配置版本
    LastConfigHash string    // 上次应用版本的哈希
    LastAppliedAt  time.Time // 上次应用时间
}
```

## 与现有代码的关系

```
runner
│
├── ProcessorRegistry          # 新增：可扩展的处理器管理器
│   └── AgentConfigProcessor  # 新增：默认 agent 配置处理器
│       └── AgentPatchDecoder  # 新增：JSON 兼容解码器
│
└── AgentRepository           # 新增抽象接口
    └── StorageAgentRepository # 新增：复用 store.Storage
            └── store.Storage.GetAgentDao()
                    └── AgentDao.ApplyHotloadPatch()  # 新增方法
                            └── models.AgentHotloadPatch # 新增：字段白名单 patch
```

所有新增代码均通过接口隔离，不破坏现有代码结构。可独立测试和替换各层实现。

## 后续规划

- [ ] 健康检查端点：暴露 `Runner.State()` 为 HTTP 接口
- [ ] 配置缓存失效 Observer：清除 viper 等本地缓存
- [ ] 指标上报 Observer：Prometheus metrics
- [ ] 熔断装饰器：gRPC 连接装饰器叠加
- [ ] 配置回滚能力：`Observer.OnConfigRollback` 扩展
- [ ] 多 Processor 并行分发（fan-out）：一个变更同时触发多个 Processor
- [ ] FileConfigProcessor：处理 refType="file" 的文件配置热更新
- [ ] BuiltinConfigProcessor：处理 refType="builtin" 的内置配置热更新