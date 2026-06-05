package vectorsearch

import (
	"strings"
	"testing"

	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/shop"
	"nova-factory-server/app/utils/store"
)

func TestBuildGoodsSearchQueryPayloads(t *testing.T) {
	categoryStore := store.NewShopCategoryStore()
	categoryStore.Set([]store.ShopCategoryData{
		&shopmodels.CategoryInfo{
			ID:           1001,
			CategoryName: "天然矿泉水",
		},
	})
	store.RegisterStore(shop.ShopCategoryStoreName, categoryStore)
	defer categoryStore.Clear()

	payloads, err := BuildGoodsSearchQueryPayloads([]string{
		"  天然矿泉水 550ml  ",
		"  69x5.3O型圈GB/T3452.1橡胶NBR  ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(payloads) != 2 {
		t.Fatalf("unexpected payload count: %#v", payloads)
	}
	if payloads[0].Original != "天然矿泉水 550ml" {
		t.Fatalf("unexpected original query: %#v", payloads[0])
	}
	if !strings.Contains(payloads[0].SearchText, "天然矿泉水") {
		t.Fatalf("expected category in search text: %#v", payloads[0])
	}
	if !strings.Contains(payloads[0].SearchText, "550ml") {
		t.Fatalf("expected spec in search text: %#v", payloads[0])
	}
	if !strings.Contains(payloads[1].SearchText, "69*5.3") {
		t.Fatalf("expected normalized spec in search text: %#v", payloads[1])
	}
}
