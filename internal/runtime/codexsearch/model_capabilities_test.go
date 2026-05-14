package codexsearch

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestParseModelCapabilitiesSearchSuffix(t *testing.T) {
	caps := ParseModelCapabilities("gpt-5-codex-search")
	if !caps.SearchEnabled {
		t.Fatalf("SearchEnabled = false, want true")
	}
	if caps.BaseModel != "gpt-5-codex" {
		t.Fatalf("BaseModel = %q, want %q", caps.BaseModel, "gpt-5-codex")
	}
	if caps.OriginalModel != "gpt-5-codex-search" {
		t.Fatalf("OriginalModel = %q", caps.OriginalModel)
	}
}

func TestParseModelCapabilitiesPlainModel(t *testing.T) {
	caps := ParseModelCapabilities("gpt-5-codex")
	if caps.SearchEnabled {
		t.Fatalf("SearchEnabled = true, want false")
	}
	if caps.BaseModel != "gpt-5-codex" {
		t.Fatalf("BaseModel = %q, want %q", caps.BaseModel, "gpt-5-codex")
	}
}

func TestParseModelCapabilitiesOnlyCompleteTrailingSearchSuffix(t *testing.T) {
	for _, model := range []string{"gpt-5-codex-searching", "gpt-5-search-codex", "search-gpt-5-codex"} {
		caps := ParseModelCapabilities(model)
		if caps.SearchEnabled {
			t.Fatalf("ParseModelCapabilities(%q).SearchEnabled = true, want false", model)
		}
		if caps.BaseModel != model {
			t.Fatalf("ParseModelCapabilities(%q).BaseModel = %q, want original", model, caps.BaseModel)
		}
	}
}

func TestInjectCachedWebSearchNoTools(t *testing.T) {
	input := []byte(`{"model":"gpt-5-codex","input":[]}`)
	out := InjectCachedWebSearch(input)
	tools := gjson.GetBytes(out, "tools")
	if len(tools.Array()) != 1 {
		t.Fatalf("tools length = %d, want 1; payload=%s", len(tools.Array()), string(out))
	}
	assertCachedWebSearchTool(t, tools.Array()[0])
}

func TestInjectCachedWebSearchKeepsFunctionTools(t *testing.T) {
	input := []byte(`{"tools":[{"type":"function","name":"lookup","parameters":{"type":"object"}}]}`)
	out := InjectCachedWebSearch(input)
	tools := gjson.GetBytes(out, "tools")
	arr := tools.Array()
	if len(arr) != 2 {
		t.Fatalf("tools length = %d, want 2; payload=%s", len(arr), string(out))
	}
	if arr[0].Get("type").String() != "function" || arr[0].Get("name").String() != "lookup" {
		t.Fatalf("function tool not preserved; payload=%s", string(out))
	}
	assertCachedWebSearchTool(t, arr[1])
}

func TestInjectCachedWebSearchKeepsImageGenerationTools(t *testing.T) {
	input := []byte(`{"tools":[{"type":"image_generation","output_format":"png"}]}`)
	out := InjectCachedWebSearch(input)
	tools := gjson.GetBytes(out, "tools")
	arr := tools.Array()
	if len(arr) != 2 {
		t.Fatalf("tools length = %d, want 2; payload=%s", len(arr), string(out))
	}
	if arr[0].Get("type").String() != "image_generation" {
		t.Fatalf("image_generation tool not preserved; payload=%s", string(out))
	}
	assertCachedWebSearchTool(t, arr[1])
}

func TestInjectCachedWebSearchDoesNotDuplicateWebSearch(t *testing.T) {
	input := []byte(`{"tools":[{"type":"web_search","search_context_size":"high"}]}`)
	out := InjectCachedWebSearch(input)
	tools := gjson.GetBytes(out, "tools")
	if len(tools.Array()) != 1 {
		t.Fatalf("tools length = %d, want 1; payload=%s", len(tools.Array()), string(out))
	}
	if got := tools.Array()[0].Get("search_context_size").String(); got != "high" {
		t.Fatalf("existing web_search not preserved; search_context_size=%q", got)
	}
}

func TestInjectCachedWebSearchDoesNotDuplicateWebSearchPreview(t *testing.T) {
	input := []byte(`{"tools":[{"type":"web_search_preview"}]}`)
	out := InjectCachedWebSearch(input)
	tools := gjson.GetBytes(out, "tools")
	if len(tools.Array()) != 1 {
		t.Fatalf("tools length = %d, want 1; payload=%s", len(tools.Array()), string(out))
	}
	if got := tools.Array()[0].Get("type").String(); got != "web_search_preview" {
		t.Fatalf("existing web_search_preview not preserved; type=%q", got)
	}
}

func assertCachedWebSearchTool(t *testing.T, tool gjson.Result) {
	t.Helper()
	if got := tool.Get("type").String(); got != "web_search" {
		t.Fatalf("tool.type = %q, want web_search", got)
	}
	if !tool.Get("external_web_access").Exists() {
		t.Fatalf("tool.external_web_access missing")
	}
	if got := tool.Get("external_web_access").Bool(); got {
		t.Fatalf("tool.external_web_access = %v, want false", got)
	}
}
