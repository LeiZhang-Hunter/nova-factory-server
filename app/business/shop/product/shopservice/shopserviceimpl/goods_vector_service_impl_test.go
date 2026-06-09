package shopserviceimpl

import (
	"strings"
	"testing"

	"nova-factory-server/app/business/shop/product/shopmodels"
)

func TestBuildGoodsEmbeddingPayloadsBySku(t *testing.T) {
	goods := &shopmodels.Goods{
		GoodsName:        "商品A",
		GoodsCode:        "goods-a",
		ShopCategoryName: "分类A",
		Description:      "商品描述",
		Skus: []*shopmodels.GoodsSku{
			{SkuID: "sku-1", SkuName: "红色", SkuCode: "red", Quantity: 5},
			{SkuID: "sku-2", SkuName: "蓝色", SkuCode: "blue", Quantity: 8},
		},
	}

	payloads := buildGoodsEmbeddingPayloads(goods)
	if len(payloads) != 2 {
		t.Fatalf("unexpected payload count: %d", len(payloads))
	}
	if payloads[0].sku == nil || payloads[0].sku.SkuID != "sku-1" {
		t.Fatalf("unexpected first sku payload")
	}
	if !strings.Contains(payloads[0].content, "SKU名称: 红色") {
		t.Fatalf("first payload missing first sku content: %s", payloads[0].content)
	}
	if strings.Contains(payloads[0].content, "SKU名称: 蓝色") {
		t.Fatalf("first payload should not contain second sku content: %s", payloads[0].content)
	}
	if payloads[1].sku == nil || payloads[1].sku.SkuID != "sku-2" {
		t.Fatalf("unexpected second sku payload")
	}
	if !strings.Contains(payloads[1].content, "SKU名称: 蓝色") {
		t.Fatalf("second payload missing second sku content: %s", payloads[1].content)
	}
	if strings.Contains(payloads[1].content, "SKU名称: 红色") {
		t.Fatalf("second payload should not contain first sku content: %s", payloads[1].content)
	}
}

func TestBuildGoodsEmbeddingPayloadsWithoutSku(t *testing.T) {
	goods := &shopmodels.Goods{
		GoodsName:   "商品B",
		GoodsCode:   "goods-b",
		Description: "商品描述B",
	}

	payloads := buildGoodsEmbeddingPayloads(goods)
	if len(payloads) != 1 {
		t.Fatalf("unexpected payload count: %d", len(payloads))
	}
	if payloads[0].sku != nil {
		t.Fatalf("unexpected sku payload")
	}
	if !strings.Contains(payloads[0].content, "商品名称: 商品B") {
		t.Fatalf("payload missing goods content: %s", payloads[0].content)
	}
}

func TestDeduplicateGoodsVectorBatchSearchRows(t *testing.T) {
	rows := []*shopmodels.GoodsVectorBatchSearchItem{
		{
			Query: "鞋子",
			Rows: []*shopmodels.GoodsVectorSearchItem{
				{GoodsDBID: 1, SkuID: 11, GoodsName: "商品1", Score: 0.9},
				{GoodsDBID: 1, SkuID: 12, GoodsName: "商品1重复", Score: 0.8},
				{GoodsDBID: 2, SkuID: 21, GoodsName: "商品2", Score: 0.7},
			},
			Total: 3,
		},
		{
			Query: "鞋子",
			Rows: []*shopmodels.GoodsVectorSearchItem{
				{GoodsDBID: 3, SkuID: 31, GoodsName: "商品3", Score: 0.6},
			},
			Total: 1,
		},
		{
			Query: "帽子",
			Rows: []*shopmodels.GoodsVectorSearchItem{
				{SkuID: 41, GoodsName: "商品4", Score: 0.5},
				{SkuID: 41, GoodsName: "商品4重复", Score: 0.4},
			},
			Total: 2,
		},
	}

	deduped := deduplicateGoodsVectorBatchSearchRows(rows)
	if len(deduped) != 2 {
		t.Fatalf("unexpected batch row count: %d", len(deduped))
	}
	if deduped[0].Query != "鞋子" || deduped[1].Query != "帽子" {
		t.Fatalf("unexpected query order: %#v", []string{deduped[0].Query, deduped[1].Query})
	}
	if len(deduped[0].Rows) != 2 || deduped[0].Total != 2 {
		t.Fatalf("unexpected first query rows: len=%d total=%d", len(deduped[0].Rows), deduped[0].Total)
	}
	if deduped[0].Rows[0].SkuID != 11 || deduped[0].Rows[1].GoodsDBID != 2 {
		t.Fatalf("unexpected first query dedup result: %#v", deduped[0].Rows)
	}
	if len(deduped[1].Rows) != 1 || deduped[1].Total != 1 || deduped[1].Rows[0].Score != 0.5 {
		t.Fatalf("unexpected second query rows: %#v", deduped[1].Rows)
	}
}
