# Gemini CLI Search 虚拟模型

本文说明本 fork 中的 Gemini CLI `-search` 虚拟模型功能。

## 功能概览

当 Gemini CLI / Code Assist 模型可用时，CPA 会通过 OpenAI-compatible `/v1/models` 接口自动暴露额外的 `-search` 变体。

示例：

```text
gemini-3-pro-preview
gemini-3-pro-preview-search
gemini-3-flash-preview
gemini-3-flash-preview-search
```

任何支持从 `/v1/models` 发现模型的客户端或网关，都可以自动看到这些虚拟模型。

## 请求行为

当客户端请求如下模型：

```text
gemini-3-pro-preview-search
```

CPA 会在发往上游前还原为真实基础模型名：

```text
gemini-3-pro-preview
```

在发送 Gemini CLI / Code Assist 上游请求前，CPA 会向上游请求体注入 Gemini 内置 `googleSearch` 工具声明。

搜索由上游 Gemini / Code Assist 服务端执行。CPA 不实现本地 search tool loop，也不会把 `googleSearch` 暴露为 OpenAI `tool_calls`。

如果请求体中已经存在 `googleSearch`，CPA 不会重复注入。用户自定义 function tools 会被保留。

## 配置项

该功能默认启用。

如需关闭自动暴露 `-search` 模型和 `googleSearch` 注入：

```yaml
disable-gemini-search-models: true
```

关闭后，CPA 保持普通 Gemini CLI 模型行为，不再把 `-search` 当作虚拟搜索后缀处理。

## 支持的模型名

解析逻辑只识别完整的 `-search` 后缀，不匹配模型名中间的 `search` 字符串。

支持示例：

```text
gemini-3-pro-preview-search
gemini-3-flash-preview-search
gemini-3-pro-preview-search(high)
gemini-3-pro-preview(high)-search
```

## 生效范围和限制

- 仅 `gemini-cli` provider 会自动生成 `-search` 变体。
- 其他 provider 不会生成 `-search` 变体。
- search 实际效果取决于上游 Gemini / Code Assist 是否接受当前模型和账号的 `googleSearch` 工具声明。
- 真实联网搜索需要部署后使用有效 Gemini CLI auth 验证。
- CPA 继续使用现有流式和非流式响应转换逻辑。

## 建议验证方式

1. 调用 `GET /v1/models`，确认普通 Gemini CLI 模型和 `-search` 变体都存在。
2. 请求 `gemini-xxx-search`，确认上游请求模型名为 `gemini-xxx`。
3. 确认上游 Gemini CLI / Code Assist 请求体包含 `googleSearch`。
4. 请求普通 `gemini-xxx`，确认不会注入 `googleSearch`。
5. 部署后使用有效 Gemini CLI auth 验证真实搜索行为。
