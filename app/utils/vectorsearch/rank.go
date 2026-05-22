package vectorsearch

import (
	"math"
	"sort"
	"strings"
)

// RankCandidate 表示进入应用层精排的候选文档。
//
// 这类候选通常已经经过第一阶段召回，例如：
// 1. 向量检索召回
// 2. BM25/全文检索召回
// 3. RRF 等融合策略初排
//
// 当前结构只保留精排所需的核心字段：
// - ID: 结果唯一标识，便于回写原始结果或做稳定排序
// - Title/Code/Category/Unit/Standard: 结构化高价值字段
// - Remark/Content: 补充语义字段
// - BaseScore: 第一阶段召回返回的原始分数
type RankCandidate struct {
	ID        int64
	Title     string
	Code      string
	Category  string
	Unit      string
	Standard  string
	Remark    string
	Content   string
	BaseScore float32
}

// RankedCandidate 表示重排后的结果索引与最终得分。
//
// 这里不直接复制完整文档内容，而是只返回：
// - Index: 对应原始 candidates 切片中的位置
// - Score: 应用层融合后的最终得分
//
// 这样可以避免重复拷贝大对象，也方便上层用 Index 回填原始业务结果。
type RankedCandidate struct {
	Index int
	Score float32
}

// RerankCandidates 对候选集进行应用层精排。
//
// 整体流程分为四步：
// 1. 归一化第一阶段召回分数，避免不同召回通道的分数尺度不一致
// 2. 结合字段匹配、编码命中、标签信号、原始排序位置计算综合分
// 3. 按综合分做稳定排序
// 4. 根据 query 类型和 top1/top2 分差动态计算阈值，过滤弱相关结果
//
// 这个函数的目标不是替代底层向量库排序，而是在“已召回候选”的基础上做更贴近业务语义的精排。
func RerankCandidates(query *ProcessedQuery, candidates []RankCandidate, limit int) []RankedCandidate {
	if len(candidates) == 0 {
		return make([]RankedCandidate, 0)
	}
	if query == nil {
		query = ProcessQuery("")
	}
	if limit <= 0 || limit > len(candidates) {
		limit = len(candidates)
	}

	// 先把底层召回分压到统一量纲，避免不同检索方式产生的分值不可直接比较。
	baseScores := normalizeBaseScores(candidates)
	ranked := make([]RankedCandidate, 0, len(candidates))
	for idx, candidate := range candidates {
		finalScore := combineCandidateScore(query, candidate, baseScores[idx], idx)
		ranked = append(ranked, RankedCandidate{
			Index: idx,
			Score: float32(finalScore),
		})
	}

	// 使用稳定排序，在分数相同的情况下保留原始顺序，减少结果抖动。
	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].Score == ranked[j].Score {
			return ranked[i].Index < ranked[j].Index
		}
		return ranked[i].Score > ranked[j].Score
	})

	threshold := calcDynamicThreshold(query, ranked)
	minKeep := 1
	if limit > 1 {
		minKeep = minInt(limit, 3)
	}

	// 至少保留前几个候选，避免阈值过严导致“明明有结果却被过滤空”的体验问题。
	filtered := make([]RankedCandidate, 0, limit)
	for idx, item := range ranked {
		if len(filtered) >= limit {
			break
		}
		if idx < minKeep || float64(item.Score) >= threshold {
			filtered = append(filtered, item)
		}
	}
	if len(filtered) == 0 {
		return ranked[:limit]
	}
	return filtered
}

// combineCandidateScore 计算单个候选的综合得分。
//
// 评分逻辑由四类信号组成：
// 1. baseScore: 底层召回的原始相似度，代表“整体相关性”
// 2. fieldScore: 各结构化字段的匹配得分，代表“业务可解释性”
// 3. tagScore: 分类/单位/规格中的最强标签信号，代表“标签命中”
// 4. rankBonus: 对靠前候选给予轻微位置奖励，减少初排信息完全丢失
//
// 默认场景下：
// - 保留 baseScore 的主导作用，但不过度依赖
// - 强化标题、编码、规格等字段的价值
//
// 编码类 query 场景下：
// - 明显提高 codeScore 权重
// - 降低正文 content 的干扰
// - 对精确码命中追加额外奖励
//
// 短 query 场景下：
// - 额外增强标题命中，因为短 query 往往更像“关键词检索”
func combineCandidateScore(query *ProcessedQuery, candidate RankCandidate, baseScore float64, rank int) float64 {
	titleScore := scoreTextField(query, candidate.Title)
	codeScore := scoreCodeField(query, candidate.Code)
	categoryScore := scoreTextField(query, candidate.Category)
	unitScore := scoreTextField(query, candidate.Unit)
	standardScore := scoreTextField(query, candidate.Standard)
	remarkScore := scoreTextField(query, candidate.Remark)
	contentScore := scoreTextField(query, candidate.Content)
	specScore := scoreSpecTerms(query, candidate.Standard, candidate.Title, candidate.Content)

	// fieldScore 体现“字段加权”思想：
	// - Title 权重最高，因为商品名通常最能代表实际意图
	// - Code 次高，因为编码/条码类查询对精确命中极为敏感
	// - Standard、Category 其次，适合规格或分类型搜索
	// - Remark、Content 权重较低，避免长文本噪声覆盖核心字段
	fieldScore := titleScore*0.28 + codeScore*0.22 + standardScore*0.14 + specScore*0.16 + categoryScore*0.09 + unitScore*0.05 + remarkScore*0.02 + contentScore*0.04
	// tagScore 只取几个标签字段中的最大值，代表“是否命中过最关键的标签维度”。
	tagScore := maxFloat(categoryScore, unitScore, standardScore, specScore)
	// rankBonus 对初排更靠前的候选给一个轻微奖励，避免精排完全无视召回顺序。
	rankBonus := 1 / float64(rank+3)

	// 默认配比下：
	// - 0.42 给原始召回分，保留向量/BM25 初排价值
	// - 0.38 给字段匹配分，强调结构化字段解释性
	// - 0.12 给标签信号，补充业务属性相关性
	// - 0.08 给位置奖励，让前序结果更稳定
	score := baseScore*0.42 + fieldScore*0.38 + tagScore*0.12 + rankBonus*0.08
	if query.IsCodeLike {
		// 编码类检索通常目标非常明确，例如条码、SKU、货号。
		// 这类场景下若仍过度依赖语义和正文，容易把“语义相关但编码不对”的结果顶上来，
		// 因此这里显著抬高 codeScore 权重，并压低正文影响。
		score = baseScore*0.24 + fieldScore*0.24 + codeScore*0.42 + tagScore*0.04 + rankBonus*0.06
		if codeScore >= 0.95 {
			// 对精确匹配或前缀高匹配再加一档奖励，尽量把“真正的目标商品”顶到最前面。
			score += 0.18
		}
		// 编码检索时正文很可能只是噪声来源，因此给予轻微惩罚。
		score -= contentScore * 0.05
	}
	if len(query.SpecTerms) > 0 {
		// 规格词命中时额外拉开与“同名不同规格”商品的差距。
		score += specScore * 0.12
		if specScore == 0 {
			// 当 query 明确带规格而候选完全没有命中规格时，给予一档惩罚，
			// 避免高 baseScore 的“同名不同规格”商品排到前面。
			score -= 0.18
		}
	}
	if query.IsShortQuery {
		// 短 query 往往缺少上下文，商品名标题的判别力更强，因此额外补一点标题权重。
		score += titleScore * 0.08
	}
	// 最终把得分限制在合理范围，避免不同规则叠加后出现异常值。
	return clampFloat(score, 0, 1.5)
}

// calcDynamicThreshold 根据当前排序结果动态估算保留阈值。
//
// 设计目标：
// - 阈值太低：会放进很多“勉强相关”的长尾结果
// - 阈值太高：会把本来可用的次优结果过度过滤
//
// 因此这里不使用固定阈值，而是根据：
// 1. top1 的绝对强度
// 2. query 是否是编码类/短 query
// 3. top1 与 top2 的分差
// 动态计算本次查询的过滤门槛。
func calcDynamicThreshold(query *ProcessedQuery, ranked []RankedCandidate) float64 {
	if len(ranked) == 0 {
		return 0
	}
	top := float64(ranked[0].Score)
	// 默认取 top1 的 58% 作为基线，兼顾召回率与精度。
	threshold := top * 0.58
	if query != nil && query.IsCodeLike {
		// 编码类检索更适合高精度，阈值抬高，减少误召回。
		threshold = top * 0.72
	} else if query != nil && query.IsShortQuery {
		// 短 query 歧义更高，适当放宽阈值，让候选不要被裁得过少。
		threshold = top * 0.52
	}
	if len(ranked) > 1 {
		second := float64(ranked[1].Score)
		gap := top - second
		switch {
		case gap < 0.05:
			// top1/top2 很接近时，说明头部结果竞争激烈，应适当放宽阈值。
			threshold -= 0.06
		case gap > 0.20:
			// top1 明显领先时，说明 query 很聚焦，可以适当收紧阈值。
			threshold += 0.05
		}
	}
	// 最终阈值限制在一个经验区间，避免极端情况过宽或过窄。
	return clampFloat(threshold, 0.18, 0.92)
}

// normalizeBaseScores 将候选的原始召回分归一化到 0~1。
//
// 原因在于不同召回通道、不同索引、不同 query 下的原始分数尺度可能不同，
// 若直接参与融合，会导致某一类分数天然占优。
//
// 这里使用 min-max 归一化：
// - 正常情况下，把最小值压到 0，最大值压到 1
// - 若所有分值几乎一样，则给一个保守的默认值，避免除零和无意义放大
func normalizeBaseScores(candidates []RankCandidate) []float64 {
	if len(candidates) == 0 {
		return nil
	}
	minScore := float64(candidates[0].BaseScore)
	maxScore := minScore
	for _, candidate := range candidates[1:] {
		score := float64(candidate.BaseScore)
		if score < minScore {
			minScore = score
		}
		if score > maxScore {
			maxScore = score
		}
	}
	result := make([]float64, 0, len(candidates))
	if math.Abs(maxScore-minScore) < 1e-9 {
		fill := 0.5
		if maxScore <= 0 {
			// 如果原始分本身也偏弱，则整体给更保守的基础分。
			fill = 0.2
		}
		for range candidates {
			result = append(result, fill)
		}
		return result
	}
	for _, candidate := range candidates {
		score := (float64(candidate.BaseScore) - minScore) / (maxScore - minScore)
		result = append(result, clampFloat(score, 0, 1))
	}
	return result
}

// scoreTextField 计算普通文本字段对 query 的匹配得分。
//
// 评分策略采用“精确命中 > 完整包含 > 关键词覆盖率 > token 覆盖率”的思路：
// - 精确等值：直接返回 1
// - 包含完整 query：给高分 0.92
// - 否则退化为按关键词和 token 的覆盖率估分
//
// 这样做的目的，是兼顾：
// - 精确匹配场景的高精度
// - 中文短语、拆词、同义词扩展后的宽召回
func scoreTextField(query *ProcessedQuery, field string) float64 {
	if query == nil {
		return 0
	}
	field = strings.ToLower(NormalizeWhitespace(field))
	if field == "" {
		return 0
	}
	if query.Normalized != "" && field == query.Normalized {
		return 1
	}
	score := 0.0
	if query.Normalized != "" && strings.Contains(field, query.Normalized) {
		score = 0.92
	}
	matched := 0
	for _, keyword := range query.Keywords {
		if keyword == "" {
			continue
		}
		if strings.Contains(field, keyword) {
			matched++
		}
	}
	if len(query.Keywords) > 0 {
		// Keywords 覆盖率更偏向“增强后的检索意图”。
		coverage := float64(matched) / float64(len(query.Keywords))
		if coverage > score {
			score = coverage
		}
	}
	if len(query.Tokens) > 0 {
		tokenMatched := 0
		for _, token := range query.Tokens {
			if token == "" {
				continue
			}
			if strings.Contains(field, token) {
				tokenMatched++
			}
		}
		// Tokens 覆盖率更偏向“原始 query 分词命中”。
		coverage := float64(tokenMatched) / float64(len(query.Tokens))
		score = maxFloat(score, coverage)
	}
	return clampFloat(score, 0, 1)
}

// scoreCodeField 计算编码类字段的匹配得分。
//
// 该函数只关注紧凑匹配，不做复杂语义判断：
// - 完全相等：1
// - 前缀命中：0.95
// - 子串包含：0.88
// - 其他：0
//
// 这是因为编码类字段的核心诉求是“精准识别”，不是“语义相近”。
func scoreCodeField(query *ProcessedQuery, field string) float64 {
	if query == nil {
		return 0
	}
	field = strings.ToLower(strings.ReplaceAll(NormalizeWhitespace(field), " ", ""))
	if field == "" {
		return 0
	}
	value := strings.ReplaceAll(query.Normalized, " ", "")
	if value == "" {
		return 0
	}
	switch {
	case field == value:
		return 1
	case strings.HasPrefix(field, value):
		return 0.95
	case strings.Contains(field, value):
		return 0.88
	default:
		return 0
	}
}

// clampFloat 将浮点数限制在给定区间内。
func clampFloat(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

// maxFloat 返回多个浮点数中的最大值。
func maxFloat(values ...float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxValue := values[0]
	for _, value := range values[1:] {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

// minInt 返回两个整数中的较小值。
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
