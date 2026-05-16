package shopdaoimpl

import (
	"testing"

	"nova-factory-server/app/business/shop/product/shopmodels"
)

func TestBuildGoodsVectorRowsWithSkus(t *testing.T) {
	goods := &shopmodels.Goods{
		ID:               1001,
		GoodsID:          "goods-1",
		GoodsName:        "商品A",
		GoodsCode:        "code-a",
		ShopCategoryName: "分类A",
		Description:      "商品描述",
	}

	items := []*shopmodels.GoodsVectorUpsertItem{
		{SkuID: 11, SkuName: "规格1", RetailPrice: 19.9, Weight: 1.2, Quantity: 8, Content: "embedding content 1", Vector: []float32{0.1, 0.2}},
		{SkuID: 22, SkuName: "规格2", RetailPrice: 29.9, Weight: 2.3, Quantity: 16, Content: "embedding content 2", Vector: []float32{0.3, 0.4}},
	}

	rows := buildGoodsVectorRows(goods, items)
	if len(rows.pks) != 2 {
		t.Fatalf("unexpected pk count: %d", len(rows.pks))
	}
	if rows.pks[0] != 11 {
		t.Fatalf("unexpected first pk: %d", rows.pks[0])
	}
	if rows.pks[1] != 22 {
		t.Fatalf("unexpected second pk: %d", rows.pks[1])
	}
	if rows.skuIDs[0] != 11 || rows.skuNames[0] != "规格1" {
		t.Fatalf("unexpected first sku fields: %d %s", rows.skuIDs[0], rows.skuNames[0])
	}
	if rows.skuIDs[1] != 22 || rows.skuNames[1] != "规格2" {
		t.Fatalf("unexpected second sku fields: %d %s", rows.skuIDs[1], rows.skuNames[1])
	}
	if rows.contents[0] != "embedding content 1" || rows.contents[1] != "embedding content 2" {
		t.Fatalf("unexpected contents: %#v", rows.contents)
	}
	if rows.retailPrices[0] != 19.9 || rows.weights[0] != 1.2 || rows.quantities[0] != 8 {
		t.Fatalf("unexpected first sku metrics: %v %v %d", rows.retailPrices[0], rows.weights[0], rows.quantities[0])
	}
	if rows.retailPrices[1] != 29.9 || rows.weights[1] != 2.3 || rows.quantities[1] != 16 {
		t.Fatalf("unexpected second sku metrics: %v %v %d", rows.retailPrices[1], rows.weights[1], rows.quantities[1])
	}
}

func TestBuildGoodsVectorRowsWithoutSku(t *testing.T) {
	goods := &shopmodels.Goods{
		ID:        1002,
		GoodsID:   "goods-2",
		GoodsName: "商品B",
	}

	items := []*shopmodels.GoodsVectorUpsertItem{
		{Content: "embedding content", Vector: []float32{0.3, 0.4}},
	}

	rows := buildGoodsVectorRows(goods, items)
	if len(rows.pks) != 1 {
		t.Fatalf("unexpected pk count: %d", len(rows.pks))
	}
	if rows.pks[0] != 0 {
		t.Fatalf("unexpected pk: %d", rows.pks[0])
	}
	if rows.skuIDs[0] != 0 || rows.skuNames[0] != "" {
		t.Fatalf("unexpected empty sku fields: %d %s", rows.skuIDs[0], rows.skuNames[0])
	}
	if rows.retailPrices[0] != 0 || rows.weights[0] != 0 || rows.quantities[0] != 0 {
		t.Fatalf("unexpected empty sku metrics: %v %v %d", rows.retailPrices[0], rows.weights[0], rows.quantities[0])
	}
}
