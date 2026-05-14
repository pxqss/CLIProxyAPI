# Search Enhancements

This document describes the provider-scoped `-search` virtual model features in this fork.

## Feature overview

This branch supports two search enhancement paths:

- Gemini CLI `-search` virtual models inject Gemini's built-in `googleSearch` tool declaration.
- Codex `-search` virtual models inject cached Codex hosted `web_search`.

## Model list behavior

When supported provider models are available, CPA automatically exposes both the base model and the corresponding `-search` variant through the OpenAI-compatible `/v1/models` endpoint. Plain base models remain available.

Example:

```text
gemini-3-pro-preview
gemini-3-pro-preview-search
gpt-5-codex
gpt-5-codex-search
```

Any client or gateway that discovers models from an OpenAI-compatible `/v1/models` endpoint can see these virtual models automatically, without manual model entry creation.

## Gemini CLI search behavior

When a client requests a model such as:

```text
gemini-xxx-search
```

CPA sends the request upstream as the real base model:

```text
gemini-xxx
```

Before sending the Gemini CLI / Code Assist upstream request, CPA injects Gemini's built-in `googleSearch` tool declaration into the upstream request body.

Search is executed by the upstream Gemini / Code Assist service. CPA does not implement a local search tool loop and does not expose `googleSearch` as OpenAI `tool_calls`.

If the request already contains `googleSearch`, CPA does not inject a duplicate declaration. User-provided function tools are preserved.

## Codex search behavior

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

Search enhancements are enabled by default.

To disable Gemini CLI `-search` model exposure and `googleSearch` injection:

```yaml
disable-gemini-search-models: true
```

To disable Codex `-search` model exposure and cached `web_search` injection:

```yaml
disable-codex-search-models: true
```

When disabled, CPA keeps normal provider model behavior and does not treat `-search` as a virtual search suffix for that provider.

## Scope and limitations

- Gemini search only applies to the `gemini-cli` provider.
- Codex search only applies to the Codex provider.
- Other providers do not automatically generate `-search` variants.
- Codex search is cached search in this first version, not live search.
- Real search behavior depends on upstream service and account capabilities and must be validated after deployment.
- In Chat Completions streaming, Codex search process events may not be displayed, but final text should keep the existing conversion behavior.
- Responses passthrough is lower risk because Responses events are mostly forwarded as-is.

## Suggested validation

1. Call `GET /v1/models` and confirm supported providers expose `-search` variants while base models remain.
2. Compare plain models with `-search` models.
3. For Gemini, ask a weather or news question using `gemini-xxx-search`.
4. For Codex, ask a weather or news question using `gpt-xxx-search` and request sources.
5. Confirm plain models do not inject search tools.
