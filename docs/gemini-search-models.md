# Gemini CLI Search Models

This document describes the Gemini CLI `-search` virtual model feature in this fork.

## Overview

When Gemini CLI / Code Assist models are available, CPA automatically exposes additional `-search` variants through the OpenAI-compatible `/v1/models` endpoint.

Example:

```text
gemini-3-pro-preview
gemini-3-pro-preview-search
gemini-3-flash-preview
gemini-3-flash-preview-search
```

Any client or gateway that discovers models from `/v1/models` can see these virtual models automatically.

## Request behavior

When a client requests a model such as:

```text
gemini-3-pro-preview-search
```

CPA sends the request upstream as the real base model:

```text
gemini-3-pro-preview
```

Before sending the Gemini CLI / Code Assist upstream request, CPA injects Gemini's built-in `googleSearch` tool declaration into the upstream request body.

Search is executed by the upstream Gemini / Code Assist service. CPA does not implement a local search tool loop and does not expose `googleSearch` as OpenAI `tool_calls`.

If the request already contains `googleSearch`, CPA does not inject a duplicate declaration. User-provided function tools are preserved.

## Configuration

This feature is enabled by default.

To disable automatic `-search` model exposure and `googleSearch` injection:

```yaml
disable-gemini-search-models: true
```

When disabled, CPA keeps normal Gemini CLI model behavior and does not treat `-search` as a virtual search suffix.

## Supported model names

The parser only recognizes a complete `-search` suffix. It does not match `search` in the middle of a model name.

Supported examples:

```text
gemini-3-pro-preview-search
gemini-3-flash-preview-search
gemini-3-pro-preview-search(high)
gemini-3-pro-preview(high)-search
```

## Scope and limitations

- Only the `gemini-cli` provider automatically generates `-search` variants.
- Other providers do not generate `-search` variants.
- Search behavior depends on whether the upstream Gemini / Code Assist service accepts the `googleSearch` tool declaration for the selected model and account.
- Real web search behavior must be validated after deployment with valid Gemini CLI authentication.
- CPA continues to use existing streaming and non-streaming response conversion logic.

## Suggested validation

1. Call `GET /v1/models` and confirm both base Gemini CLI models and `-search` variants are present.
2. Request `gemini-xxx-search` and confirm the upstream request model is `gemini-xxx`.
3. Confirm the upstream Gemini CLI / Code Assist request body contains `googleSearch`.
4. Request plain `gemini-xxx` and confirm `googleSearch` is not injected.
5. Validate real search behavior after deployment with valid Gemini CLI authentication.
