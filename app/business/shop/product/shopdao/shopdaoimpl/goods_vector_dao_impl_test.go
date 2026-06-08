package shopdaoimpl

import (
	"encoding/json"
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
		IsOnSale:         1,
	}

	items := []*shopmodels.GoodsVectorUpsertItem{
		{SkuID: 11, SkuName: "规格1", RetailPrice: 19.9, Weight: 1.2, Quantity: 8, Content: "embedding content 1", Vector: []float32{0.1, 0.2}},
		{SkuID: 22, SkuName: "规格2", RetailPrice: 29.9, Weight: 2.3, Quantity: 16, Content: "embedding content 2", Vector: []float32{0.3, 0.4}},
	}

	rows, err := buildGoodsVectorRows(goods, items, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows.skuIDs) != 2 {
		t.Fatalf("unexpected sku count: %d", len(rows.skuIDs))
	}
	if rows.skuIDs[0] != 11 {
		t.Fatalf("unexpected first sku id: %d", rows.skuIDs[0])
	}
	if rows.skuIDs[1] != 22 {
		t.Fatalf("unexpected second sku id: %d", rows.skuIDs[1])
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
	if !rows.saleFlags[0] || !rows.saleFlags[1] {
		t.Fatalf("unexpected sale flags: %#v", rows.saleFlags)
	}
}

func TestBuildGoodsVectorRowsWithoutSku(t *testing.T) {
	goods := &shopmodels.Goods{
		ID:        1002,
		GoodsID:   "goods-2",
		GoodsName: "商品B",
		IsOnSale:  0,
	}

	items := []*shopmodels.GoodsVectorUpsertItem{
		{Content: "embedding content", Vector: []float32{0.3, 0.4}},
	}

	rows, err := buildGoodsVectorRows(goods, items, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows.skuIDs) != 1 {
		t.Fatalf("unexpected sku count: %d", len(rows.skuIDs))
	}
	if rows.skuIDs[0] != 0 {
		t.Fatalf("unexpected sku id: %d", rows.skuIDs[0])
	}
	if rows.skuIDs[0] != 0 || rows.skuNames[0] != "" {
		t.Fatalf("unexpected empty sku fields: %d %s", rows.skuIDs[0], rows.skuNames[0])
	}
	if rows.retailPrices[0] != 0 || rows.weights[0] != 0 || rows.quantities[0] != 0 {
		t.Fatalf("unexpected empty sku metrics: %v %v %d", rows.retailPrices[0], rows.weights[0], rows.quantities[0])
	}
	if rows.saleFlags[0] {
		t.Fatalf("unexpected sale flag: %#v", rows.saleFlags)
	}
}

func TestBuildGoodsVectorRowsWithExplicitMetadata(t *testing.T) {
	goods := &shopmodels.Goods{
		ID:               1003,
		GoodsID:          "goods-3",
		GoodsName:        "商品C",
		ShopCategoryName: "默认分类",
	}

	items := []*shopmodels.GoodsVectorUpsertItem{
		{
			SkuID:    33,
			Content:  "测试内容",
			Vector:   []float32{0.5, 0.6},
			Metadata: map[string]any{"category": "显式分类", "spec": "69*5.3", "shopId": 1001},
		},
	}

	rows, err := buildGoodsVectorRows(goods, items, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows.metadatas) != 1 {
		t.Fatalf("unexpected metadata count: %d", len(rows.metadatas))
	}

	metadata := make(map[string]any)
	if err = json.Unmarshal(rows.metadatas[0], &metadata); err != nil {
		t.Fatalf("unmarshal metadata fail: %v", err)
	}
	if metadata["category"] != "显式分类" {
		t.Fatalf("unexpected category metadata: %#v", metadata)
	}
	if metadata["spec"] != "69*5.3" {
		t.Fatalf("unexpected spec metadata: %#v", metadata)
	}
	if metadata["shopId"] != "1001" {
		t.Fatalf("unexpected shopId metadata: %#v", metadata)
	}
}
