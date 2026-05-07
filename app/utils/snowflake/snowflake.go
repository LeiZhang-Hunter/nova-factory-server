package snowflake

import (
	"github.com/spf13/viper"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var sfNode *sf.Node

func init() {
	var id int64 = 1
	var st time.Time
	start_time := viper.GetString("start_time")
	if start_time == "" {
		start_time = "2022-08-08"
	}
	st, err := time.Parse("2006-01-02", start_time)
	if err != nil {
		panic(err)
	}
	sf.Epoch = st.UnixNano() / 1000000
	sfNode, err = sf.NewNode(id)
	if err != nil {
		panic(err)
	}
}

func GenID() int64 {
	return sfNode.Generate().Int64()
}
