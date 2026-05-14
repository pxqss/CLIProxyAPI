# Search Enhancements / 搜索增强

本文说明本 fork 中按 provider 生效的 `-search` 虚拟模型能力。

## 功能概览

本分支支持两类搜索增强：

- Gemini CLI `-search` 虚拟模型：注入 Gemini 内置 `googleSearch` 工具声明。
- Codex `-search` 虚拟模型：注入 cached Codex hosted `web_search`。

## 模型列表行为

当支持的 provider 模型可用时，CPA 会通过 OpenAI-compatible `/v1/models` 接口自动暴露基础模型和对应的 `-search` 变体。普通基础模型仍然保留。

示例：

```text
gemini-3-pro-preview
gemini-3-pro-preview-search
gpt-5-codex
gpt-5-codex-search
```

任何支持从 OpenAI-compatible `/v1/models` 发现模型的客户端或网关，都可以自动看到这些虚拟模型，无需手动新增模型条目。

## Gemini CLI search 行为

当客户端请求如下模型：

```text
gemini-xxx-search
```

CPA 会在发往上游前还原为真实基础模型名：

```text
gemini-xxx
```

在发送 Gemini CLI / Code Assist 上游请求前，CPA 会向上游请求体注入 Gemini 内置 `googleSearch` 工具声明。

搜索由上游 Gemini / Code Assist 服务端执行。CPA 不实现本地 search tool loop，也不会把 `googleSearch` 暴露为 OpenAI `tool_calls`。

如果请求体中已经存在 `googleSearch`，CPA 不会重复注入。用户自定义 function tools 会被保留。

## Codex search 行为

当客户端请求如下模型：

```text
gpt-xxx-search
```

CPA 会在发往上游前还原为真实基础模型名：

```text
gpt-xxx
```

在发送 Codex / Responses 上游请求前，CPA 会注入 cached hosted web search 工具声明：

```json
{"type":"web_search","external_web_access":false}
```

这是 Codex CLI 默认 cached search 对应的 hosted tool 形式。搜索由上游 Codex / Responses 服务端执行。CPA 不实现本地 search tool loop，也不会把 `web_search` 伪装成客户端 function tool。

已有 function tools、`image_generation` tools 和其他 builtin tools 会被保留。如果请求体中已经存在 `web_search` 或 `web_search_preview`，CPA 不会重复注入。

## 配置项

搜索增强默认启用。

如需关闭 Gemini CLI `-search` 模型自动暴露和 `googleSearch` 注入：

```yaml
disable-gemini-search-models: true
```

如需关闭 Codex `-search` 模型自动暴露和 cached `web_search` 注入：

```yaml
disable-codex-search-models: true
```

关闭后，CPA 会保持对应 provider 的普通模型行为，不再把 `-search` 当作该 provider 的虚拟搜索后缀处理。

## 生效范围和限制

- Gemini search 仅 `gemini-cli` provider 生效。
- Codex search 仅 Codex provider 生效。
- 其他 provider 不会自动生成 `-search` 变体。
- Codex 第一版是 cached search，不是 live search。
- 真实搜索效果取决于上游服务和账号能力，需要部署后验证。
- Chat Completions 流式下 Codex 搜索过程事件可能不显示，但最终文本应保持现有转换逻辑。
- Responses passthrough 风险较低，因为 Responses 事件大多按现有逻辑透传。

## 建议验证方式

1. 调用 `GET /v1/models`，确认支持的 provider 会暴露 `-search` 变体，且基础模型仍然保留。
2. 对比普通模型和 `-search` 模型。
3. Gemini 使用 `gemini-xxx-search` 测试天气或新闻问题。
4. Codex 使用 `gpt-xxx-search` 测试天气或新闻问题，并要求返回来源。
5. 确认普通模型不会注入 search tool。
