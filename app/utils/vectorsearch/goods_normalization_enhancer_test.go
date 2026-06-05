package vectorsearch

import (
	"testing"

	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/shop"
	"nova-factory-server/app/utils/store"
)

func TestGoodsNormalizationEnhancerInitAndNormalization(t *testing.T) {
	categoryStore := store.NewShopCategoryStore()
	categoryStore.Set([]store.ShopCategoryData{
		&shopmodels.CategoryInfo{
			ID:           1001,
			CategoryName: "天然矿泉水",
		},
	})
	store.RegisterStore(shop.ShopCategoryStoreName, categoryStore)
	defer categoryStore.Clear()

	enhancer := NewGoodsNormalizationEnhancer()
	if enhancer == nil {
		t.Fatalf("expected enhancer")
	}
	if err := enhancer.Init(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := enhancer.Normalization("  21x5.3O型圈GB/T3452.1橡胶NBR  ")
	if got != "天然矿泉水 69*5.3" {
		t.Fatalf("unexpected normalized value: %q", got)
	}
}

func TestGoodsNormalizationEnhancerNormalizationWithEmptyInput(t *testing.T) {
	enhancer := NewGoodsNormalizationEnhancer()
	if err := enhancer.Init(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := enhancer.Normalization("   ")
	if got != "" {
		t.Fatalf("expected empty normalized value, got %q", got)
	}
}
