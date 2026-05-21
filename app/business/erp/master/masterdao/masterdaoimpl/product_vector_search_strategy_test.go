package masterdaoimpl

import (
	"testing"

	"nova-factory-server/app/business/erp/master/mastermodels"
)

func TestBuildProductVectorSearchPlans_CodeIntent(t *testing.T) {
	data, err := buildProductVectorSearchPlans(&mastermodels.ProductVectorBatchSearchReq{
		Queries: []string{"SP-1001"},
		Limit:   10,
	}, [][]float32{{0.1, 0.2}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("unexpected plans length: %d", len(data))
	}
	if data[0].Intent != "code" {
		t.Fatalf("expected code intent, got %s", data[0].Intent)
	}
	if data[0].ExactCodeText == "" {
		t.Fatalf("expected exact code text")
	}
}

func TestBuildProductVectorSearchPlans_SpecIntent(t *testing.T) {
	data, err := buildProductVectorSearchPlans(&mastermodels.ProductVectorBatchSearchReq{
		Queries: []string{"矿泉水 550ml"},
		Limit:   10,
	}, [][]float32{{0.1, 0.2}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("unexpected plans length: %d", len(data))
	}
	if data[0].Intent != "spec" {
		t.Fatalf("expected spec intent, got %s", data[0].Intent)
	}
}

func TestBuildProductVectorSearchPlans_SemanticIntent(t *testing.T) {
	data, err := buildProductVectorSearchPlans(&mastermodels.ProductVectorBatchSearchReq{
		Queries: []string{"饮用水"},
		Limit:   10,
	}, [][]float32{{0.1, 0.2}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("unexpected plans length: %d", len(data))
	}
	if data[0].Intent != "semantic" {
		t.Fatalf("expected semantic intent, got %s", data[0].Intent)
	}
}

func TestBuildProductVectorSearchPlans_CategoryIntent(t *testing.T) {
	data, err := buildProductVectorSearchPlans(&mastermodels.ProductVectorBatchSearchReq{
		Queries: []string{"饮料"},
		Limit:   10,
	}, [][]float32{{0.1, 0.2}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("unexpected plans length: %d", len(data))
	}
	if data[0].Intent != "category" {
		t.Fatalf("expected category intent, got %s", data[0].Intent)
	}
}
