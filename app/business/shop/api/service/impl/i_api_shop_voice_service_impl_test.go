//go:build ai

package impl

import "testing"

func TestCleanTextForVoiceSynthesis(t *testing.T) {
	input := "# 标题\n\n- **第一项**：[链接](https://example.com)\n> 引用内容\n\n```go\nfmt.Println(\"hi\")\n```\n| 名称 | 值 |\n| --- | --- |\n| A | 1 |"

	got := cleanTextForVoiceSynthesis(input)
	want := `标题 第一项：链接 引用内容 fmt.Println"hi" 名称 ， 值 ， A ， 1`
	if got != want {
		t.Fatalf("cleanTextForVoiceSynthesis() = %q, want %q", got, want)
	}
}
