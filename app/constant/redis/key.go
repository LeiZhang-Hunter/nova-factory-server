package redis

import "fmt"

const AGENT_HEADETBEAT_CACHE = "agent_heartbeat_time_"

const AgentProcessCacheKey = "agent_process_key_%s_%d"
const AIAgentAliveCacheKey = "ai:agent:alive:%d"
const AIGatewayAliveCacheKey = "ai:gateway:alive:%d"

const CameraInfoCacheKey = "camera_info_%d"

func MakeCacheKey(key string, cid string, objectId uint64) string {
	return fmt.Sprintf(key, cid, objectId)
}

func MakeAIAgentAliveCacheKey(agentID int64) string {
	return fmt.Sprintf(AIAgentAliveCacheKey, agentID)
}

func MakeAIGatewayAliveCacheKey(gatewayID int64) string {
	return fmt.Sprintf(AIGatewayAliveCacheKey, gatewayID)
}

// IntegrationLoginCacheKeyPattern erp集成系统登陆数据
const IntegrationLoginCacheKeyPattern = "erp:integration:login:%s"
