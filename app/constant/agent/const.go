package agent

var (
	AgentStateAll     = 0
	AgentStateOnline  = 2
	AgentStateOffline = 1
)

var CHECK_ONLINE_DURATION int = 600
var CACHE_LIVE_TIME int = 1800
