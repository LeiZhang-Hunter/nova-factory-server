package goods

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"

	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/shop"
	"nova-factory-server/app/utils/store"
)

func TestMetadataExtractorInitAndExtract(t *testing.T) {
	categoryStore := store.NewShopCategoryStore()
	categoryStore.Set([]store.ShopCategoryData{
		&shopmodels.CategoryInfo{
			ID:           1001,
			CategoryName: "天然矿泉水",
		},
	})
	store.RegisterStore(shop.ShopCategoryStoreName, categoryStore)
	defer categoryStore.Clear()

	viper.Set("vectorsearch.normalization.goods.normalize.interceptors", []map[string]any{
		{
			"name": "trim_whitespace",
			"type": "whitespace",
		},
		{
			"name": "detect_category",
			"type": "product_type",
		},
		{
			"name": "replace_lower_x",
			"type": "replace",
			"old":  "x",
			"new":  "*",
		},
		{
			"name": "replace_upper_x",
			"type": "replace",
			"old":  "X",
			"new":  "*",
		},
	})
	viper.Set("vectorsearch.normalization.goods.metadata.spec", []map[string]any{
		{
			"name":        "extract_spec",
			"type":        "regex",
			"pattern":     `(?i)^.*?((?:d)?\d+(?:\.\d+)?\*\d+(?:\.\d+)?\*\d+(?:\.\d+)?|(?:d)?\d+(?:\.\d+)?\*\d+(?:\.\d+)?(?:[-\s]+(?:d)?\d+(?:\.\d+)?\*\d+(?:\.\d+)?)?).*$`,
			"replacement": "$1",
		},
	})
	viper.Set("vectorsearch.normalization.goods.metadata.category.interceptors", []map[string]any{
		{
			"name": "detect_category",
			"type": "product_type",
		},
	})
	t.Cleanup(func() {
		viper.Set("vectorsearch.normalization.goods.normalize", nil)
		viper.Set("vectorsearch.normalization.goods.metadata", nil)
	})

	extractor := NewMetadataExtractor()
	if extractor == nil {
		t.Fatalf("expected extractor")
	}
	if err := extractor.Init(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result, err := extractor.Extract("  天然矿泉水 69x5.3O型圈GB/T3452.1橡胶NBR  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Value != "天然矿泉水 69*5.3O型圈GB/T3452.1橡胶NBR" {
		t.Fatalf("unexpected normalized value: %q", result.Value)
	}
	if len(result.Categories) != 1 || result.Categories[0].Name != "天然矿泉水" || result.Categories[0].ID != 1001 {
		t.Fatalf("unexpected categories: %#v", result.Categories)
	}
	if len(result.Metadata["category"]) != 1 || result.Metadata["category"][0] != "天然矿泉水" {
		t.Fatalf("unexpected category metadata: %#v", result.Metadata)
	}
	if len(result.Metadata["spec"]) == 0 || result.Metadata["spec"][0] != "69*5.3" {
		t.Fatalf("unexpected spec metadata: %#v", result.Metadata)
	}
	fmt.Println(11)
	return
}
