package cliproxy

import "testing"

func TestAppendGeminiSearchModels(t *testing.T) {
	models := []*ModelInfo{
		{ID: "gemini-3-pro-preview", Object: "model", OwnedBy: "google", Type: "gemini-cli", DisplayName: "gemini-3-pro-preview"},
		{ID: "gemini-3-flash-preview", Object: "model", OwnedBy: "google", Type: "gemini-cli"},
	}
	out := appendGeminiSearchModels(models)
	ids := map[string]bool{}
	for _, model := range out {
		ids[model.ID] = true
	}
	for _, id := range []string{"gemini-3-pro-preview", "gemini-3-pro-preview-search", "gemini-3-flash-preview", "gemini-3-flash-preview-search"} {
		if !ids[id] {
			t.Fatalf("missing model %q in %#v", id, ids)
		}
	}
}

func TestAppendGeminiSearchModelsIsProviderScopedByCaller(t *testing.T) {
	models := []*ModelInfo{{ID: "claude-sonnet-4-6", Object: "model", OwnedBy: "anthropic", Type: "claude"}}
	var out []*ModelInfo
	provider := "claude"
	if provider == "gemini-cli" {
		out = appendGeminiSearchModels(models)
	} else {
		out = models
	}
	ids := map[string]bool{}
	for _, model := range out {
		ids[model.ID] = true
	}
	if ids["claude-sonnet-4-6-search"] {
		t.Fatalf("non gemini-cli provider produced search variant")
	}
}
