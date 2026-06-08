package normalization

import (
	"testing"

	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/shop"
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/store"
	"nova-factory-server/app/utils/vectorsearch/normalization/api"
	"nova-factory-server/app/utils/vectorsearch/normalization/lowercase"
	"nova-factory-server/app/utils/vectorsearch/normalization/regex"
	replacepkg "nova-factory-server/app/utils/vectorsearch/normalization/replace"
	"nova-factory-server/app/utils/vectorsearch/normalization/shopcategory"
	"nova-factory-server/app/utils/vectorsearch/normalization/whitespace"
)

func TestNewPipelineWithConfiguredSteps(t *testing.T) {
	regexProperties, err := cfg.Pack(regex.Config{
		Pattern:     `[-_/]`,
		Replacement: "",
	})
	if err != nil {
		t.Fatalf("pack regex config failed: %v", err)
	}

	pipeline := NewPipeline(api.Config{
		Interceptors: []*api.InterceptorConfig{
			{Type: whitespace.Type},
			{Type: lowercase.Type},
			{Type: regex.Type, Properties: regexProperties},
		},
	})

	result, err := pipeline.Normalize("  AB_C-12/3  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Value != "abc123" {
		t.Fatalf("unexpected normalized value: %q", result.Value)
	}
	if len(result.Matches) != 1 || result.Matches[0].Kind != "regex_replace" {
		t.Fatalf("unexpected matches: %#v", result.Matches)
	}
}

func TestNewPipelineWithCategoryStep(t *testing.T) {
	categoryStore := store.NewShopCategoryStore()
	categoryStore.Set([]store.ShopCategoryData{
		&shopmodels.CategoryInfo{
			ID:           1001,
			CategoryName: "饮用水",
		},
	})
	store.RegisterStore(shop.ShopCategoryStoreName, categoryStore)
	defer categoryStore.Clear()

	pipeline := NewPipeline(api.Config{
		Interceptors: []*api.InterceptorConfig{
			{Type: whitespace.Type},
			{Type: shopcategory.Type},
		},
	})

	result, err := pipeline.Normalize("  饮用水  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Categories) != 1 || result.Categories[0].Name != "饮用水" || result.Categories[0].ID != 1001 {
		t.Fatalf("unexpected categories: %#v", result.Categories)
	}
	if len(result.Metadata["category"]) != 1 || result.Metadata["category"][0] != "饮用水" {
		t.Fatalf("unexpected metadata: %#v", result.Metadata)
	}
	if len(result.Metadata["category_id"]) != 1 || result.Metadata["category_id"][0] != "1001" {
		t.Fatalf("unexpected category id metadata: %#v", result.Metadata)
	}
}

func TestNewPipelineWithCategoryStepMatchFromCache(t *testing.T) {
	categoryStore := store.NewShopCategoryStore()
	categoryStore.Set([]store.ShopCategoryData{
		&shopmodels.CategoryInfo{
			ID:           2001,
			CategoryName: "饮用水",
			Children: []*shopmodels.CategoryInfo{
				{
					ID:           2002,
					CategoryName: "天然矿泉水",
				},
			},
		},
	})
	store.RegisterStore(shop.ShopCategoryStoreName, categoryStore)
	defer categoryStore.Clear()

	pipeline := NewPipeline(api.Config{
		Interceptors: []*api.InterceptorConfig{
			{Type: whitespace.Type},
			{Type: shopcategory.Type},
		},
	})

	result, err := pipeline.Normalize("  天然矿泉水 550ml ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Categories) == 0 || result.Categories[0].Name != "天然矿泉水" || result.Categories[0].ID != 2002 {
		t.Fatalf("unexpected categories: %#v", result.Categories)
	}
	if len(result.Metadata["category_id"]) == 0 || result.Metadata["category_id"][0] != "2002" {
		t.Fatalf("unexpected category id metadata: %#v", result.Metadata)
	}
}

func TestNewPipelineWithReplaceStep(t *testing.T) {
	replaceProperties, err := cfg.Pack(replacepkg.Config{
		Old: "x",
		New: "*",
	})
	if err != nil {
		t.Fatalf("pack replace config failed: %v", err)
	}
	regexProperties, err := cfg.Pack(regex.Config{
		Pattern:     `(?i)^.*?((?:d)?\d+(?:\.\d+)?\*\d+(?:\.\d+)?\*\d+(?:\.\d+)?|(?:d)?\d+(?:\.\d+)?\*\d+(?:\.\d+)?(?:[-\s]+(?:d)?\d+(?:\.\d+)?\*\d+(?:\.\d+)?)?).*$`,
		Replacement: "$1",
	})
	if err != nil {
		t.Fatalf("pack regex config failed: %v", err)
	}

	pipeline := NewPipeline(api.Config{
		Interceptors: []*api.InterceptorConfig{
			{Type: whitespace.Type},
			{Type: replacepkg.Type, Properties: replaceProperties},
			{Type: regex.Type, Properties: regexProperties},
		},
	})

	cases := []struct {
		name           string
		input          string
		want           string
		wantMatchCount int
	}{
		{
			name:           "three_dimensions",
			input:          " 5型试样PTFE棕+二硫化钼-4.5*4.5*320 ",
			want:           "4.5*4.5*320",
			wantMatchCount: 1,
		},
		{
			name:           "replace_then_extract",
			input:          " 69x5.3O型圈GB/T3452.1橡胶NBR  ",
			want:           "69*5.3",
			wantMatchCount: 2,
		},
		{
			name:           "double_spec_range",
			input:          "ZHIDE-黑色OR90N-d19.5*1.5-D22.5*1.5",
			want:           "d19.5*1.5-D22.5*1.5",
			wantMatchCount: 1,
		},
		{
			name:           "double_spec_range",
			input:          "PTFE绿-113*145*120",
			want:           "113*145*120",
			wantMatchCount: 1,
		},
		{
			name:           "double_spec_range",
			input:          "黑色磨砂氟胶O型圈-d101.32*1.78 D104.88*1.78",
			want:           "d101.32*1.78 D104.88*1.78",
			wantMatchCount: 1,
		},
	}

	for _, tc := range cases {
		result, err := pipeline.Normalize(tc.input)
		if err != nil {
			t.Fatalf("%s: expected no error, got %v", tc.name, err)
		}
		if result.Value != tc.want {
			t.Fatalf("%s: unexpected normalized value: %q", tc.name, result.Value)
		}
		if len(result.Matches) != tc.wantMatchCount {
			t.Fatalf("%s: unexpected matches: %#v", tc.name, result.Matches)
		}
		if tc.wantMatchCount == 2 {
			if result.Matches[0].Kind != "replace" || result.Matches[1].Kind != "regex_replace" {
				t.Fatalf("%s: unexpected matches: %#v", tc.name, result.Matches)
			}
			continue
		}
		if result.Matches[0].Kind != "regex_replace" {
			t.Fatalf("%s: unexpected matches: %#v", tc.name, result.Matches)
		}
	}
	return
}
