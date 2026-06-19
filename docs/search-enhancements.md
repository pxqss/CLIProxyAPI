# Codex Search Enhancements

This document describes the Codex-only `-search` virtual model feature in this fork.

## Feature Overview

This branch supports Codex `-search` virtual models. CPA injects the cached Codex hosted `web_search` tool into upstream Codex / Responses requests.

Gemini CLI / GCLI search is no longer supported. The upstream project removed the Gemini CLI OAuth path, so this fork no longer exposes Gemini `-search` virtual models.

## Model List Behavior

When Codex models are available, CPA automatically exposes both the base model and the corresponding `-search` variant through the OpenAI-compatible `/v1/models` endpoint. Plain base models remain available.

Example:

```text
gpt-5-codex
gpt-5-codex-search
```

Any client or gateway that discovers models from an OpenAI-compatible `/v1/models` endpoint can see these virtual models automatically, without manual model entry creation.

## Codex Search Behavior

When a client requests a model such as:

```text
gpt-xxx-search
```

CPA sends the request upstream as the real base model:

```text
gpt-xxx
```

Before sending the Codex / Responses upstream request, CPA injects the cached hosted web search tool declaration:

```json
{"type":"web_search","external_web_access":false}
```

This is the hosted tool form corresponding to Codex CLI's default cached search. Search is executed by the upstream Codex / Responses service. CPA does not implement a local search tool loop and does not disguise `web_search` as a client function tool.

Existing function tools, `image_generation` tools, and other built-in tools are preserved. If the request already contains `web_search` or `web_search_preview`, CPA does not inject a duplicate declaration.

## Configuration

Codex search enhancements are enabled by default.

To disable Codex `-search` model exposure and cached `web_search` injection:

```yaml
disable-codex-search-models: true
```

When disabled, CPA keeps normal Codex model behavior and does not treat `-search` as a virtual search suffix for Codex.

## Scope And Limitations

- Codex search only applies to the Codex provider.
- Gemini CLI / GCLI OAuth search is not supported.
- Gemini, Claude, Qwen, OpenCode, Antigravity, and other providers do not automatically generate Codex `-search` variants.
- Codex search is cached search in this first version, not live search.
- Real search behavior depends on upstream service and account capabilities and must be validated after deployment.
- In Chat Completions streaming, Codex search process events may not be displayed, but final text should keep the existing conversion behavior.
- Responses passthrough is lower risk because Responses events are mostly forwarded as-is.

## Suggested Validation

1. Call `GET /v1/models` and confirm Codex exposes `-search` variants while base models remain.
2. Compare a plain Codex model with its `-search` variant.
3. Ask a weather or news question using `gpt-xxx-search` and request sources.
4. Confirm plain Codex models do not inject `web_search`.
