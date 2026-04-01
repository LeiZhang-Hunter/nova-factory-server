package baizeContext

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ParamInt64(c *gin.Context, key string) int64 {
	i, err := strconv.ParseInt(c.Param(key), 10, 64)
	if err != nil {
		return 0
	}
	return i
}
func ParamInt64Array(c *gin.Context, key string) []int64 {
	split := strings.Split(c.Param(key), ",")
	is := make([]int64, 0, len(split))
	for _, s := range split {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil
		}
		is = append(is, i)
	}
	return is
}
func ParamStringArray(c *gin.Context, key string) []string {
	return strings.Split(c.Param(key), ",")
}

func QueryInt64(c *gin.Context, key string) int64 {
	i, err := strconv.ParseInt(c.Query(key), 10, 64)
	if err != nil {
		return 0
	}
	return i
}
func QueryInt64Array(c *gin.Context, key string) []int64 {
	split := strings.Split(c.Query(key), ",")
	is := make([]int64, 0, len(split))
	for _, s := range split {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil
		}
		is = append(is, i)
	}
	return is

}

func QueryUint64(c *gin.Context, key string) uint64 {
	i, err := strconv.ParseUint(c.Query(key), 10, 64)
	if err != nil {
		return 0
	}
	return i
}
