package uuid

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"github.com/satori/go.uuid"
)

func MakeUuid() string {
	h := md5.New()
	h.Write(uuid.NewV4().Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

func HashId() uint64 {
	has := xxhash.New()
	has.WriteString(MakeUuid())
	h4 := has.Sum64()
	fmt.Println(h4) //输出753694413698530628
	return (h4)
}
