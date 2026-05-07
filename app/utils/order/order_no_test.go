package order

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

func TestGenerateOrderNo(t *testing.T) {
	viper.Set("start_time", "2022-08-08")
	first := GenerateOrderNo()
	fmt.Println(first)
	second := GenerateOrderNo()

	if first[:3] != "ORD" {
		t.Fatalf("unexpected order no prefix: %s", first)
	}
	if first == second {
		t.Fatalf("generated duplicate order no: %s", first)
	}
	if len(first) <= len("ORD") {
		t.Fatalf("unexpected order no value: %s", first)
	}
}

func TestGenerateOrderNoWithPrefix(t *testing.T) {
	orderNo := GenerateOrderNoWithPrefix("wx")
	if len(orderNo) <= len("wx") {
		t.Fatalf("unexpected order no value: %s", orderNo)
	}
	if orderNo[:2] != "wx" {
		t.Fatalf("unexpected order no prefix: %s", orderNo)
	}
}
