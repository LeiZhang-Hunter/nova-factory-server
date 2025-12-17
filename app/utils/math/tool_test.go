package math

import (
	"fmt"
	"testing"
)

func TestRoundFloat(t *testing.T) {
	v := RoundFloat(4.98822, 2)
	fmt.Println(v)
}
