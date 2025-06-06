package redis

import "fmt"

const AGENT_HEADETBEAT_CACHE = "agent_heartbeat_time_"

const AgentProcessCacheKey = "agent_process_key_%s_%d"

func MakeCacheKey(key string, cid string, objectId uint64) string {
	return fmt.Sprintf(key, cid, objectId)
}
