package executor

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v7/internal/config"
	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v7/sdk/cliproxy/auth"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v7/sdk/cliproxy/executor"
	sdktranslator "github.com/router-for-me/CLIProxyAPI/v7/sdk/translator"
	"github.com/tidwall/gjson"
)

func TestCodexExecutorExecuteSearchModelInjectsCachedWebSearch(t *testing.T) {
	gotBody := executeCodexRequestAndCaptureBody(t, "gpt-5-codex-search")
	if got := gjson.GetBytes(gotBody, "model").String(); got != "gpt-5-codex" {
		t.Fatalf("model = %q, want gpt-5-codex; body=%s", got, string(gotBody))
	}
	assertCodexCachedWebSearchInjected(t, gotBody)
}

func TestCodexExecutorExecutePlainModelDoesNotInjectWebSearch(t *testing.T) {
	gotBody := executeCodexRequestAndCaptureBody(t, "gpt-5-codex")
	if got := gjson.GetBytes(gotBody, "model").String(); got != "gpt-5-codex" {
		t.Fatalf("model = %q, want gpt-5-codex; body=%s", got, string(gotBody))
	}
	assertCodexWebSearchAbsent(t, gotBody)
}

func executeCodexRequestAndCaptureBody(t *testing.T, model string) []byte {
	t.Helper()
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		gotBody = body
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_1\",\"object\":\"response\",\"created_at\":0,\"status\":\"completed\",\"background\":false,\"error\":null,\"output\":[]}}\n\n"))
	}))
	defer server.Close()

	executor := NewCodexExecutor(&config.Config{SDKConfig: config.SDKConfig{DisableImageGeneration: config.DisableImageGenerationAll}})
	auth := &cliproxyauth.Auth{Attributes: map[string]string{
		"base_url": server.URL,
		"api_key":  "test",
	}}

	_, err := executor.Execute(context.Background(), auth, cliproxyexecutor.Request{
		Model:   model,
		Payload: []byte(`{"model":"` + model + `","input":"hello"}`),
	}, cliproxyexecutor.Options{
		SourceFormat: sdktranslator.FromString("openai-response"),
		Stream:       false,
	})
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	return gotBody
}

func assertCodexCachedWebSearchInjected(t *testing.T, body []byte) {
	t.Helper()
	tools := gjson.GetBytes(body, "tools")
	if !tools.IsArray() {
		t.Fatalf("tools missing or not array; body=%s", string(body))
	}
	count := 0
	for _, tool := range tools.Array() {
		if tool.Get("type").String() != "web_search" {
			continue
		}
		count++
		if !tool.Get("external_web_access").Exists() {
			t.Fatalf("web_search external_web_access missing; body=%s", string(body))
		}
		if got := tool.Get("external_web_access").Bool(); got {
			t.Fatalf("external_web_access = %v, want false; body=%s", got, string(body))
		}
	}
	if count != 1 {
		t.Fatalf("web_search tool count = %d, want 1; body=%s", count, string(body))
	}
}

func assertCodexWebSearchAbsent(t *testing.T, body []byte) {
	t.Helper()
	tools := gjson.GetBytes(body, "tools")
	if !tools.IsArray() {
		return
	}
	for _, tool := range tools.Array() {
		if tool.Get("type").String() == "web_search" || tool.Get("type").String() == "web_search_preview" {
			t.Fatalf("unexpected web_search tool; body=%s", string(body))
		}
	}
}
