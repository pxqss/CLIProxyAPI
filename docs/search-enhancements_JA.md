# Codex Search Enhancements / Codex 検索拡張

このドキュメントでは、この fork における Codex provider 専用の `-search` 仮想モデル機能について説明します。

## 機能概要

このブランチは Codex `-search` 仮想モデルをサポートします。CPA は上流の Codex / Responses リクエストに cached Codex hosted `web_search` ツールを注入します。

Gemini CLI / GCLI search はサポートされません。上流プロジェクトで Gemini CLI OAuth 経路が削除されたため、この fork でも Gemini `-search` 仮想モデルは公開しません。

## モデル一覧の動作

Codex モデルが利用可能な場合、CPA は OpenAI-compatible `/v1/models` エンドポイントを通じて、ベースモデルと対応する `-search` バリアントを自動的に公開します。通常のベースモデルも引き続き利用できます。

例:

```text
gpt-5-codex
gpt-5-codex-search
```

OpenAI-compatible `/v1/models` からモデルを発見できる任意のクライアントやゲートウェイで、これらの仮想モデルを自動的に確認でき、手動でモデルを追加する必要はありません。

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

Codex 検索拡張はデフォルトで有効です。

Codex `-search` モデルの自動公開と cached `web_search` 注入を無効化するには、次を設定します。

```yaml
disable-codex-search-models: true
```

無効化すると、CPA は Codex の通常モデル動作を維持し、`-search` を Codex の仮想検索サフィックスとして扱いません。

## 適用範囲と制限

- Codex search は Codex provider のみに適用されます。
- Gemini CLI / GCLI OAuth search はサポートされません。
- Gemini、Claude、Qwen、OpenCode、Antigravity、およびその他の provider では Codex `-search` バリアントを自動生成しません。
- Codex search の初回版は cached search であり、live search ではありません。
- 実際の検索挙動は上流サービスとアカウント能力に依存するため、デプロイ後の検証が必要です。
- Chat Completions streaming では Codex の検索過程イベントが表示されない場合がありますが、最終テキストは既存の変換ロジックを維持する想定です。
- Responses passthrough は、Responses イベントの多くが既存ロジックでそのまま透過されるため、比較的リスクが低いです。

## 推奨検証手順

1. `GET /v1/models` を呼び出し、Codex が `-search` バリアントを公開し、ベースモデルも残っていることを確認します。
2. 通常の Codex モデルと対応する `-search` モデルを比較します。
3. `gpt-xxx-search` を使って天気やニュースの質問をテストし、出典を要求します。
4. 通常の Codex モデルで `web_search` が注入されないことを確認します。
