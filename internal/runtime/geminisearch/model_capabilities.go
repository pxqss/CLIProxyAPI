package geminisearch

import (
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const SearchSuffix = "-search"

type Capabilities struct {
	OriginalModel string
	BaseModel     string
	SearchEnabled bool
}

// ParseModelCapabilities only treats a complete -search suffix as virtual search.
func ParseModelCapabilities(model string) Capabilities {
	caps := Capabilities{OriginalModel: model, BaseModel: model}
	trimmed := strings.TrimSpace(model)
	if trimmed == "" {
		return caps
	}

	if base, ok := stripSearchBeforeThinking(trimmed); ok {
		caps.BaseModel = base
		caps.SearchEnabled = true
		return caps
	}
	if base, ok := stripSearchAfterThinking(trimmed); ok {
		caps.BaseModel = base
		caps.SearchEnabled = true
		return caps
	}
	return caps
}

func AppendSearchVariant(model string) string {
	trimmed := strings.TrimSpace(model)
	if trimmed == "" {
		return model
	}
	if strings.HasSuffix(trimmed, ")") {
		if open := strings.LastIndex(trimmed, "("); open > 0 {
			name := trimmed[:open]
			return name + SearchSuffix + trimmed[open:]
		}
	}
	return trimmed + SearchSuffix
}

func InjectGoogleSearch(payload []byte) []byte {
	if len(payload) == 0 || hasGoogleSearch(payload) {
		return payload
	}
	tools := gjson.GetBytes(payload, "request.tools")
	if tools.Exists() && tools.IsArray() {
		out, err := sjson.SetRawBytes(payload, "request.tools.-1", []byte(`{"googleSearch":{}}`))
		if err == nil {
			return out
		}
		return payload
	}
	out, err := sjson.SetRawBytes(payload, "request.tools", []byte(`[{"googleSearch":{}}]`))
	if err == nil {
		return out
	}
	return payload
}

func hasGoogleSearch(payload []byte) bool {
	tools := gjson.GetBytes(payload, "request.tools")
	if !tools.Exists() || !tools.IsArray() {
		return false
	}
	found := false
	tools.ForEach(func(_, tool gjson.Result) bool {
		if tool.Get("googleSearch").Exists() || tool.Get("google_search").Exists() {
			found = true
			return false
		}
		return true
	})
	return found
}

func stripSearchBeforeThinking(model string) (string, bool) {
	if !strings.HasSuffix(model, ")") {
		return "", false
	}
	open := strings.LastIndex(model, "(")
	if open <= 0 {
		return "", false
	}
	name := model[:open]
	if !strings.HasSuffix(name, SearchSuffix) {
		return "", false
	}
	baseName := strings.TrimSuffix(name, SearchSuffix)
	if baseName == "" {
		return "", false
	}
	return baseName + model[open:], true
}

func stripSearchAfterThinking(model string) (string, bool) {
	if !strings.HasSuffix(model, SearchSuffix) {
		return "", false
	}
	base := strings.TrimSuffix(model, SearchSuffix)
	if base == "" {
		return "", false
	}
	return base, true
}
