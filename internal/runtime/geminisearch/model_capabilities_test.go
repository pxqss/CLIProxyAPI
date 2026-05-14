package geminisearch

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestParseModelCapabilitiesSearchSuffix(t *testing.T) {
	caps := ParseModelCapabilities("gemini-3-pro-preview-search")
	if !caps.SearchEnabled {
		t.Fatalf("SearchEnabled = false, want true")
	}
	if caps.BaseModel != "gemini-3-pro-preview" {
		t.Fatalf("BaseModel = %q, want %q", caps.BaseModel, "gemini-3-pro-preview")
	}
	if caps.OriginalModel != "gemini-3-pro-preview-search" {
		t.Fatalf("OriginalModel = %q", caps.OriginalModel)
	}
}

func TestParseModelCapabilitiesPlainModel(t *testing.T) {
	caps := ParseModelCapabilities("gemini-3-pro-preview")
	if caps.SearchEnabled {
		t.Fatalf("SearchEnabled = true, want false")
	}
	if caps.BaseModel != "gemini-3-pro-preview" {
		t.Fatalf("BaseModel = %q, want %q", caps.BaseModel, "gemini-3-pro-preview")
	}
}

func TestParseModelCapabilitiesThinkingCombinations(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "search before thinking", in: "gemini-3-pro-preview-search(high)", want: "gemini-3-pro-preview(high)"},
		{name: "search after thinking", in: "gemini-3-pro-preview(high)-search", want: "gemini-3-pro-preview(high)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caps := ParseModelCapabilities(tt.in)
			if !caps.SearchEnabled {
				t.Fatalf("SearchEnabled = false, want true")
			}
			if caps.BaseModel != tt.want {
				t.Fatalf("BaseModel = %q, want %q", caps.BaseModel, tt.want)
			}
		})
	}
}

func TestInjectGoogleSearchDoesNotDuplicate(t *testing.T) {
	input := []byte(`{"request":{"tools":[{"googleSearch":{}}]}}`)
	out := InjectGoogleSearch(input)
	tools := gjson.GetBytes(out, "request.tools")
	if len(tools.Array()) != 1 {
		t.Fatalf("tools length = %d, want 1; payload=%s", len(tools.Array()), string(out))
	}
}

func TestInjectGoogleSearchKeepsFunctionTools(t *testing.T) {
	input := []byte(`{"request":{"tools":[{"functionDeclarations":[{"name":"lookup","parametersJsonSchema":{"type":"object"}}]}]}}`)
	out := InjectGoogleSearch(input)
	tools := gjson.GetBytes(out, "request.tools")
	if len(tools.Array()) != 2 {
		t.Fatalf("tools length = %d, want 2; payload=%s", len(tools.Array()), string(out))
	}
	if !tools.Array()[0].Get("functionDeclarations.0.name").Exists() {
		t.Fatalf("functionDeclarations not preserved; payload=%s", string(out))
	}
	if !tools.Array()[1].Get("googleSearch").Exists() {
		t.Fatalf("googleSearch not appended; payload=%s", string(out))
	}
}
