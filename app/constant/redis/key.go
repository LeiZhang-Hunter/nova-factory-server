package redis

import "fmt"

const AGENT_HEADETBEAT_CACHE = "agent_heartbeat_time_"

func MakeCacheKey(key string, cid string, objectId uint64) string {
	return fmt.Sprintf(key, cid, objectId)
}
