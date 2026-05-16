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
