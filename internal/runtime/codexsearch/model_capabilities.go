package codexsearch

import (
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const SearchSuffix = "-search"

var cachedWebSearchToolJSON = []byte(`{"type":"web_search","external_web_access":false}`)

type Capabilities struct {
	OriginalModel string
	BaseModel     string
	SearchEnabled bool
}

func ParseModelCapabilities(model string) Capabilities {
	caps := Capabilities{OriginalModel: model, BaseModel: model}
	trimmed := strings.TrimSpace(model)
	if trimmed == "" || !strings.HasSuffix(trimmed, SearchSuffix) {
		return caps
	}

	base := strings.TrimSuffix(trimmed, SearchSuffix)
	if base == "" {
		return caps
	}

	caps.BaseModel = base
	caps.SearchEnabled = true
	return caps
}

func AppendSearchVariant(model string) string {
	trimmed := strings.TrimSpace(model)
	if trimmed == "" {
		return model
	}
	if strings.HasSuffix(trimmed, SearchSuffix) {
		return trimmed
	}
	return trimmed + SearchSuffix
}

func InjectCachedWebSearch(payload []byte) []byte {
	if len(payload) == 0 || hasWebSearch(payload) {
		return payload
	}

	tools := gjson.GetBytes(payload, "tools")
	if tools.Exists() && tools.IsArray() {
		out, err := sjson.SetRawBytes(payload, "tools.-1", cachedWebSearchToolJSON)
		if err == nil {
			return out
		}
		return payload
	}

	out, err := sjson.SetRawBytes(payload, "tools", []byte(`[{"type":"web_search","external_web_access":false}]`))
	if err == nil {
		return out
	}
	return payload
}

func hasWebSearch(payload []byte) bool {
	tools := gjson.GetBytes(payload, "tools")
	if !tools.Exists() || !tools.IsArray() {
		return false
	}

	found := false
	tools.ForEach(func(_, tool gjson.Result) bool {
		switch tool.Get("type").String() {
		case "web_search", "web_search_preview", "web_search_preview_2025_03_11":
			found = true
			return false
		default:
			return true
		}
	})
	return found
}
