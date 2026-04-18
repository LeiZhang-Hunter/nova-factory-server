package order

import "testing"

func TestGenerateOrderNo(t *testing.T) {
	first := GenerateOrderNo()
	second := GenerateOrderNo()

	if len(first) != 26 {
		t.Fatalf("unexpected order no length: %d, value: %s", len(first), first)
	}
	if first[:3] != "ORD" {
		t.Fatalf("unexpected order no prefix: %s", first)
	}
	if first == second {
		t.Fatalf("generated duplicate order no: %s", first)
	}
}
