package gatewayserviceimpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"nova-factory-server/app/utils/store/key"
	uuid2 "nova-factory-server/app/utils/uuid"
	"regexp"
	"strings"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AgentConfigKeyServiceImpl API Key Service 实现。
// 负责 API Key 的增删改查、生成、校验，以及 MCP 允许工具列表的管理。
type AgentConfigKeyServiceImpl struct {
	dao gatewaydao.IAgentConfigKeyDao // 底层数据访问接口，操作 ai_agent_config_key 表
}

// NewAgentConfigKeyService 创建 AgentConfigKeyServiceImpl 实例。
// 初始化后会自动注册到全局 key store，供后续 API Key 鉴权使用。
// 参数:
//   - dao: AgentConfigKey DAO 接口实现
//
// 返回:
//   - gatewayservice.IAgentConfigKeyService: 服务接口实例
func NewAgentConfigKeyService(dao gatewaydao.IAgentConfigKeyDao) gatewayservice.IAgentConfigKeyService {
	keyConfig := &AgentConfigKeyServiceImpl{dao: dao}
	// 将当前实例注册到全局 key store，用于通过 API Key 反查用户信息
	key.RegisterStore(keyConfig)
	return keyConfig
}

// Create 新增一个 API Key 记录。
// 新增前会校验 Key 格式（sk- 前缀 + 32 位 hex），并检查 Key 是否已存在。
// 参数:
//   - c: gin 上下文，用于传递用户身份等元数据
//   - req: 包含 Key 信息的保存参数
//
// 返回:
//   - *gatewaymodels.AgentConfigKey: 新创建的 API Key 记录
//   - error: 校验失败、Key 已存在或数据库写入异常时返回错误
func (a *AgentConfigKeyServiceImpl) Create(c *gin.Context, req *gatewaymodels.AgentConfigKeyUpsert) (*gatewaymodels.AgentConfigKey, error) {
	// 1. 校验 Key 格式（非空 + 正则匹配 sk-xxxxxxxx 格式）
	if err := a.validateUpsert(req); err != nil {
		return nil, err
	}
	// 2. 检查 Key 是否已被占用
	existing, err := a.dao.GetByKey(c, req.Key)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("该 API Key 已存在")
	}
	// 3. 写入数据库
	return a.dao.Create(c, req)
}

// Update 修改已有 API Key。
// 要求传入有效的 ID，且新 Key 格式需通过校验。
// 参数:
//   - c: gin 上下文
//   - req: 包含 ID 和新 Key 值的保存参数
//
// 返回:
//   - *gatewaymodels.AgentConfigKey: 更新后的完整记录
//   - error: ID 为空、校验失败或数据库更新异常时返回错误
func (a *AgentConfigKeyServiceImpl) Update(c *gin.Context, req *gatewaymodels.AgentConfigKeyUpsert) (*gatewaymodels.AgentConfigKey, error) {
	if req.ID == 0 {
		return nil, errors.New("id不能为空")
	}
	// 校验新 Key 格式
	if err := a.validateUpsert(req); err != nil {
		return nil, err
	}
	return a.dao.Update(c, req)
}

// DeleteByIDs 批量删除 API Key 记录。
// 参数:
//   - c: gin 上下文，dao 层会根据 create_by 做权限隔离
//   - ids: 待删除的 API Key ID 列表
//
// 返回:
//   - error: 删除失败时返回错误
func (a *AgentConfigKeyServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return a.dao.DeleteByIDs(c, ids)
}

// GetByID 根据 ID 获取单个 API Key 记录。
// 参数:
//   - c: gin 上下文
//   - id: API Key 记录主键
//
// 返回:
//   - *gatewaymodels.AgentConfigKey: 查询到的记录，不存在时返回 nil
//   - error: 数据库异常时返回错误
func (a *AgentConfigKeyServiceImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.AgentConfigKey, error) {
	return a.dao.GetByID(c, id)
}

// List 分页查询 API Key 列表。
// 支持按 Key 模糊搜索，默认每页 20 条，按 id 倒序排列。
// 参数:
//   - c: gin 上下文
//   - req: 查询条件，包含 Key 关键字、Page、Size；为 nil 时使用默认分页
//
// 返回:
//   - *gatewaymodels.AgentConfigKeyListData: 包含 rows 和 total 的分页结果
//   - error: 数据库查询失败时返回错误
func (a *AgentConfigKeyServiceImpl) List(c *gin.Context, req *gatewaymodels.AgentConfigKeyQuery) (*gatewaymodels.AgentConfigKeyListData, error) {
	if req == nil {
		req = new(gatewaymodels.AgentConfigKeyQuery)
	}
	return a.dao.List(c, req)
}

// Generate 生成一个随机的 API Key 字符串。
// 生成规则：取 UUID v4，去掉连字符后计算 MD5，最终格式为 "sk-" + 32位小写 hex。
// 注意：该方法只生成字符串，不写入数据库，调用方需自行调用 Create 持久化。
// 返回:
//   - string: 格式形如 "sk-abc123..." 的 API Key
func (a *AgentConfigKeyServiceImpl) Generate() string {
	// 生成 UUID v4 并去掉连字符
	uuidStr := strings.ReplaceAll(uuid.NewString(), "-", "")
	// 对 UUID 字符串做 MD5，拼上 sk- 前缀
	key := fmt.Sprintf("sk-%s", uuid2.MakeMd5([]byte(uuidStr)))
	return key
}

// validateUpsert 校验新增/修改时 Key 字段的合法性。
// 校验规则：
//  1. 不能为空
//  2. 必须符合 "sk-" 前缀 + 32 位小写十六进制字符的格式
//
// 参数:
//   - req: 包含 Key 字段的请求体
//
// 返回:
//   - error: 校验不通过时返回具体错误描述
func (a *AgentConfigKeyServiceImpl) validateUpsert(req *gatewaymodels.AgentConfigKeyUpsert) error {
	if strings.TrimSpace(req.Key) == "" {
		return errors.New("密钥不能为空")
	}
	// 正则：sk- 开头，后接 32 位小写十六进制字符 [0-9a-f]
	if !regexp.MustCompile(`^sk-[0-9a-f]{32}$`).MatchString(strings.TrimSpace(req.Key)) {
		return errors.New("密钥格式不正确，需为 sk- 前缀加 32 位小写 hex")
	}
	return nil
}

// GetUserId 根据 API Key 字符串反查对应的用户 ID（create_by）。
// 主要用于外部调用通过 Key 鉴权时，获取该 Key 所属的用户。
// 参数:
//   - key: API Key 字符串
//
// 返回:
//   - int64: 该 Key 创建者的用户 ID，查询失败或 Key 不存在时返回 0
func (a *AgentConfigKeyServiceImpl) GetUserId(key string) int64 {
	info, err := a.dao.GetByKey(&gin.Context{}, key)
	if err != nil {
		return 0
	}
	return info.CreateBy
}

// normalizeAllowMcpTools 对传入的工具名列表做规范化处理。
// 处理步骤：
//  1. 去除每个工具名的首尾空白字符
//  2. 过滤掉空字符串
//  3. 去除重复的工具名（保留首次出现的顺序）
//
// 参数:
//   - tools: 原始工具名列表，可能包含空白、重复项
//
// 返回:
//   - []string: 去重去空格后的工具名列表；如果入参为空则返回 nil
func normalizeAllowMcpTools(tools []string) []string {
	if len(tools) == 0 {
		return nil
	}

	// 预分配容量，减少切片扩容
	normalized := make([]string, 0, len(tools))
	// 用空结构体 map 记录已出现的工具名，空结构体不占内存
	seen := make(map[string]struct{}, len(tools))
	for _, tool := range tools {
		// 去除首尾空白
		tool = strings.TrimSpace(tool)
		// 过滤空白工具名
		if tool == "" {
			continue
		}
		// 去重
		if _, ok := seen[tool]; ok {
			continue
		}
		seen[tool] = struct{}{}
		normalized = append(normalized, tool)
	}
	return normalized
}

// SetAllowMcpTools 设置某个 API Key 允许使用的 MCP 工具列表。
// 业务含义：在 MCP 服务探测完成后，用户可选择允许哪些工具对外暴露，
// 该列表以 JSON 数组形式存入 ai_agent_config_key 表的 allow_mcp_server_tools 字段。
// 参数:
//   - c: gin 上下文
//   - req: 包含 API Key ID 和工具名列表的请求参数
//
// 返回:
//   - *gatewaymodels.AgentConfigKey: 更新后的完整 API Key 记录
//   - error: 参数非法、Key 不存在或数据库更新失败时返回错误
func (a *AgentConfigKeyServiceImpl) SetAllowMcpTools(c *gin.Context, req *gatewaymodels.AgentConfigKeyToolUpsert) (*gatewaymodels.AgentConfigKey, error) {
	// 1. 参数校验
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	if req.ID == 0 {
		return nil, errors.New("id不能为空")
	}

	// 2. 确认目标 API Key 记录存在
	current, err := a.dao.GetByID(c, req.ID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, errors.New("API Key不存在")
	}

	// 3. 规范化工具列表：去空格、去空值、去重
	normalizedTools := normalizeAllowMcpTools(req.Tools)
	// 若规范化后无有效工具，将字段置为空字符串（表示不限制或未配置）
	if len(normalizedTools) == 0 {
		return a.dao.UpdateAllowMcpTools(c, req.ID, "")
	}

	// 4. 将工具名列表序列化为 JSON 数组格式存入数据库
	body, err := json.Marshal(normalizedTools)
	if err != nil {
		return nil, fmt.Errorf("允许使用的MCP工具配置编码失败: %w", err)
	}
	return a.dao.UpdateAllowMcpTools(c, req.ID, string(body))
}
