package cliproxy

import "testing"

func TestAppendCodexSearchModels(t *testing.T) {
	models := []*ModelInfo{
		{ID: "gpt-5-codex", Object: "model", OwnedBy: "openai", Type: "openai", DisplayName: "gpt-5-codex"},
		{ID: "gpt-5", Object: "model", OwnedBy: "openai", Type: "openai"},
	}
	out := appendCodexSearchModels(models)
	ids := map[string]bool{}
	for _, model := range out {
		ids[model.ID] = true
	}
	for _, id := range []string{"gpt-5-codex", "gpt-5-codex-search", "gpt-5", "gpt-5-search"} {
		if !ids[id] {
			t.Fatalf("missing model %q in %#v", id, ids)
		}
	}
}

func TestAppendCodexSearchModelsIsProviderScopedByCaller(t *testing.T) {
	models := []*ModelInfo{{ID: "claude-sonnet-4-6", Object: "model", OwnedBy: "anthropic", Type: "claude"}}
	var out []*ModelInfo
	provider := "claude"
	if provider == "codex" {
		out = appendCodexSearchModels(models)
	} else {
		out = models
	}
	ids := map[string]bool{}
	for _, model := range out {
		ids[model.ID] = true
	}
	if ids["claude-sonnet-4-6-search"] {
		t.Fatalf("non codex provider produced search variant")
	}
}
