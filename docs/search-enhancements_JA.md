# Search Enhancements / 検索拡張

このドキュメントでは、この fork における provider-scoped な `-search` 仮想モデル機能について説明します。

## 機能概要

このブランチでは、次の 2 種類の検索拡張をサポートします。

- Gemini CLI `-search` 仮想モデル: Gemini 組み込みの `googleSearch` ツール宣言を注入します。
- Codex `-search` 仮想モデル: cached Codex hosted `web_search` を注入します。

## モデル一覧の動作

サポート対象 provider のモデルが利用可能な場合、CPA は OpenAI-compatible `/v1/models` エンドポイントを通じて、ベースモデルと対応する `-search` バリアントの両方を自動的に公開します。通常のベースモデルも引き続き利用できます。

例:

```text
gemini-3-pro-preview
gemini-3-pro-preview-search
gpt-5-codex
gpt-5-codex-search
```

OpenAI-compatible `/v1/models` からモデルを発見できる任意のクライアントやゲートウェイで、これらの仮想モデルを自動的に確認でき、手動でモデルを追加する必要はありません。

## Gemini CLI search の動作

クライアントが次のようなモデルをリクエストした場合:

```text
gemini-xxx-search
```

CPA は上流へ送信する前に、実際のベースモデル名へ戻します。

```text
gemini-xxx
```

Gemini CLI / Code Assist の上流リクエストを送信する前に、CPA は上流リクエスト本文へ Gemini 組み込みの `googleSearch` ツール宣言を注入します。

検索は上流の Gemini / Code Assist サービス側で実行されます。CPA はローカル search tool loop を実装せず、`googleSearch` を OpenAI `tool_calls` として公開しません。

リクエスト本文に既に `googleSearch` が含まれている場合、CPA は重複して注入しません。ユーザー定義の function tools は保持されます。

## Codex search の動作

クライアントが次のようなモデルをリクエストした場合:

```text
gpt-xxx-search
```

CPA は上流へ送信する前に、実際のベースモデル名へ戻します。

```text
gpt-xxx
```

Codex / Responses の上流リクエストを送信する前に、CPA は cached hosted web search ツール宣言を注入します。

```json
{"type":"web_search","external_web_access":false}
```

これは Codex CLI のデフォルト cached search に対応する hosted tool 形式です。検索は上流の Codex / Responses サービス側で実行されます。CPA はローカル search tool loop を実装せず、`web_search` をクライアント function tool に偽装しません。

既存の function tools、`image_generation` tools、およびその他の built-in tools は保持されます。リクエスト本文に既に `web_search` または `web_search_preview` が含まれている場合、CPA は重複して注入しません。

## 設定

検索拡張はデフォルトで有効です。

Gemini CLI `-search` モデルの自動公開と `googleSearch` 注入を無効化するには、次を設定します。

```yaml
disable-gemini-search-models: true
```

Codex `-search` モデルの自動公開と cached `web_search` 注入を無効化するには、次を設定します。

```yaml
disable-codex-search-models: true
```

無効化すると、CPA は該当 provider の通常モデル動作を維持し、`-search` をその provider の仮想検索サフィックスとして扱いません。

## 適用範囲と制限

- Gemini search は `gemini-cli` provider のみに適用されます。
- Codex search は Codex provider のみに適用されます。
- その他の provider では `-search` バリアントを自動生成しません。
- Codex search の初回版は cached search であり、live search ではありません。
- 実際の検索挙動は上流サービスとアカウント能力に依存するため、デプロイ後の検証が必要です。
- Chat Completions streaming では Codex の検索過程イベントが表示されない場合がありますが、最終テキストは既存の変換ロジックを維持する想定です。
- Responses passthrough は、Responses イベントの多くが既存ロジックでほぼそのまま転送されるため、比較的リスクが低いです。

## 推奨検証手順

1. `GET /v1/models` を呼び出し、サポート対象 provider が `-search` バリアントを公開し、ベースモデルも残っていることを確認します。
2. 通常モデルと `-search` モデルを比較します。
3. Gemini では `gemini-xxx-search` を使って天気やニュースの質問をテストします。
4. Codex では `gpt-xxx-search` を使って天気やニュースの質問をテストし、出典を要求します。
5. 通常モデルでは search tool が注入されないことを確認します。
