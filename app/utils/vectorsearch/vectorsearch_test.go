package vectorsearch

import "testing"

func TestProcessQuery(t *testing.T) {
	query := ProcessQuery("  矿泉水  ")
	if query == nil {
		t.Fatalf("expected processed query")
	}
	if query.Original != "矿泉水" {
		t.Fatalf("unexpected original query: %s", query.Original)
	}
	if len(query.ExpandedTokens) == 0 {
		t.Fatalf("expected expanded tokens")
	}
	found := false
	for _, token := range query.ExpandedTokens {
		if token == "纯净水" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected synonym token in expanded tokens: %#v", query.ExpandedTokens)
	}
	if query.EmbeddingText == query.Original {
		t.Fatalf("expected embedding text to include enhancement terms")
	}
}

func TestProcessQueryWithSpec(t *testing.T) {
	query := ProcessQuery("农夫山泉 0.55L")
	if query == nil {
		t.Fatalf("expected processed query")
	}
	if len(query.SpecTerms) == 0 {
		t.Fatalf("expected spec terms")
	}
	found := false
	for _, spec := range query.SpecTerms {
		if spec == "550ml" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected normalized spec term 550ml, got %#v", query.SpecTerms)
	}
}

func TestProcessQueryWithCategory(t *testing.T) {
	query := ProcessQuery("饮料")
	if query == nil {
		t.Fatalf("expected processed query")
	}
	if len(query.CategoryTerms) == 0 {
		t.Fatalf("expected category terms")
	}
	if query.CategoryTerms[0] != "饮料" {
		t.Fatalf("expected category term 饮料, got %#v", query.CategoryTerms)
	}
}

func TestRerankCandidates(t *testing.T) {
	query := ProcessQuery("SP-1001")
	ranked := RerankCandidates(query, []RankCandidate{
		{
			ID:        1,
			Title:     "矿泉水 550ml",
			Code:      "SP-1001",
			Content:   "天然饮用水",
			BaseScore: 0.62,
		},
		{
			ID:        2,
			Title:     "矿泉水 330ml",
			Code:      "AB-2002",
			Content:   "条码相近商品",
			BaseScore: 0.91,
		},
	}, 2)
	if len(ranked) == 0 {
		t.Fatalf("expected reranked results")
	}
	if ranked[0].Index != 0 {
		t.Fatalf("expected exact code match to rank first, got %#v", ranked)
	}
}

func TestRerankCandidatesWithSpec(t *testing.T) {
	query := ProcessQuery("矿泉水 550ml")
	ranked := RerankCandidates(query, []RankCandidate{
		{
			ID:        1,
			Title:     "矿泉水 330ml",
			Standard:  "330ml",
			Content:   "饮用水",
			BaseScore: 0.95,
		},
		{
			ID:        2,
			Title:     "矿泉水 550ml",
			Standard:  "0.55L",
			Content:   "饮用水",
			BaseScore: 0.80,
		},
	}, 2)
	if len(ranked) == 0 {
		t.Fatalf("expected reranked results")
	}
	if ranked[0].Index != 1 {
		t.Fatalf("expected spec matched item to rank first, got %#v", ranked)
	}
}

func TestRerankInventoryPrefersInStockWhenRelevanceIsClose(t *testing.T) {
	query := &ProcessedQuery{Normalized: "矿泉水", Keywords: []string{"矿泉水"}, Tokens: []string{"矿泉水"}}
	ranked := RerankCandidates(query, []RankCandidate{
		{
			ID:             1,
			Title:          "矿泉水 550ml",
			BaseScore:      0.8,
			Quantity:       0,
			InventoryKnown: true,
		},
		{
			ID:             2,
			Title:          "矿泉水 550ml",
			BaseScore:      0.8,
			Quantity:       10,
			InventoryKnown: true,
		},
	}, 2)
	if len(ranked) == 0 {
		t.Fatalf("expected reranked results")
	}
	if ranked[0].Index != 1 {
		t.Fatalf("expected in-stock candidate to rank first when relevance is close, got %#v", ranked)
	}
}

func TestRerankInventoryKeepsStrongRelevanceAheadOfWeakInStock(t *testing.T) {
	query := &ProcessedQuery{Normalized: "矿泉水 550ml", Keywords: []string{"矿泉水", "550ml"}, Tokens: []string{"矿泉水", "550ml"}}
	ranked := RerankCandidates(query, []RankCandidate{
		{
			ID:             1,
			Title:          "矿泉水 550ml",
			Standard:       "550ml",
			BaseScore:      0.8,
			Quantity:       0,
			InventoryKnown: true,
		},
		{
			ID:             2,
			Title:          "抽纸",
			Content:        "日用品",
			BaseScore:      0.8,
			Quantity:       100,
			InventoryKnown: true,
		},
	}, 2)
	if len(ranked) == 0 {
		t.Fatalf("expected reranked results")
	}
	if ranked[0].Index != 0 {
		t.Fatalf("expected stronger relevance to stay first, got %#v", ranked)
	}
}

func TestRerankInventoryKeepsExactCodeAheadOfInStockMismatch(t *testing.T) {
	query := &ProcessedQuery{Normalized: "sp-1001", Keywords: []string{"sp-1001"}, Tokens: []string{"sp-1001"}, IsCodeLike: true}
	ranked := RerankCandidates(query, []RankCandidate{
		{
			ID:             1,
			Title:          "矿泉水 550ml",
			Code:           "SP-1001",
			BaseScore:      0.8,
			Quantity:       0,
			InventoryKnown: true,
		},
		{
			ID:             2,
			Title:          "矿泉水 330ml",
			Code:           "AB-2002",
			BaseScore:      0.8,
			Quantity:       20,
			InventoryKnown: true,
		},
	}, 2)
	if len(ranked) == 0 {
		t.Fatalf("expected reranked results")
	}
	if ranked[0].Index != 0 {
		t.Fatalf("expected exact code match to stay first, got %#v", ranked)
	}
}
