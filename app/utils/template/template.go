package template

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"text/template"
)

func escape(s string) string {
	return strconv.Quote(s)
}
func pruneEscape(s string) string {
	raw := strconv.Quote(s)
	raw = strings.TrimPrefix(raw, `"`)
	raw = strings.TrimSuffix(raw, `"`)
	return raw
}

func MakeTemplate(temp string) (*template.Template, error) {
	t, err := template.New("alertTemplate").Funcs(template.FuncMap{
		"escape":      escape,
		"pruneEscape": pruneEscape,
	}).Parse(temp)
	if err != nil {
		return nil, errors.WithMessagef(err, "fail to generate template %s", temp)
	}

	return t, nil
}
