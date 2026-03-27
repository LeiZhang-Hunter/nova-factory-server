package redis

import "fmt"

const AGENT_HEADETBEAT_CACHE = "agent_heartbeat_time_"

const AgentProcessCacheKey = "agent_process_key_%s_%d"

const CameraInfoCacheKey = "camera_info_%d"

func MakeCacheKey(key string, cid string, objectId uint64) string {
	return fmt.Sprintf(key, cid, objectId)
}

// IntegrationLoginCacheKeyPattern erp集成系统登陆数据
const IntegrationLoginCacheKeyPattern = "erp:integration:login:%s"
