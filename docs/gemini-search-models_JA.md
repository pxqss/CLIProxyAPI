# Gemini CLI Search Models

このドキュメントでは、この fork における Gemini CLI `-search` 仮想モデル機能について説明します。

## 概要

Gemini CLI / Code Assist モデルが利用可能な場合、CPA は OpenAI-compatible `/v1/models` エンドポイントを通じて追加の `-search` バリアントを自動的に公開します。

例:

```text
gemini-3-pro-preview
gemini-3-pro-preview-search
gemini-3-flash-preview
gemini-3-flash-preview-search
```

`/v1/models` からモデルを発見できる任意のクライアントやゲートウェイで、これらの仮想モデルを自動的に確認できます。

## リクエスト動作

クライアントが次のようなモデルをリクエストした場合:

```text
gemini-3-pro-preview-search
```

CPA は上流へ送信する前に、実際のベースモデル名へ戻します。

```text
gemini-3-pro-preview
```

Gemini CLI / Code Assist の上流リクエストを送信する前に、CPA は上流リクエスト本文へ Gemini 組み込みの `googleSearch` ツール宣言を注入します。

検索は上流の Gemini / Code Assist サービス側で実行されます。CPA はローカル search tool loop を実装せず、`googleSearch` を OpenAI `tool_calls` として公開しません。

リクエスト本文に既に `googleSearch` が含まれている場合、CPA は重複して注入しません。ユーザー定義の function tools は保持されます。

## 設定

この機能はデフォルトで有効です。

`-search` モデルの自動公開と `googleSearch` 注入を無効化するには、次を設定します。

```yaml
disable-gemini-search-models: true
```

無効化すると、CPA は通常の Gemini CLI モデル動作を維持し、`-search` を仮想検索サフィックスとして扱いません。

## サポートされるモデル名

パーサーは完全な `-search` サフィックスのみを認識します。モデル名の途中にある `search` 文字列には一致しません。

サポート例:

```text
gemini-3-pro-preview-search
gemini-3-flash-preview-search
gemini-3-pro-preview-search(high)
gemini-3-pro-preview(high)-search
```

## 適用範囲と制限

- `-search` バリアントを自動生成するのは `gemini-cli` provider のみです。
- 他の provider では `-search` バリアントを生成しません。
- search の実際の挙動は、上流の Gemini / Code Assist が選択されたモデルとアカウントに対して `googleSearch` ツール宣言を受け入れるかどうかに依存します。
- 実際の Web 検索は、デプロイ後に有効な Gemini CLI auth で検証する必要があります。
- CPA は既存のストリーミングおよび非ストリーミングのレスポンス変換ロジックを引き続き使用します。

## 推奨検証手順

1. `GET /v1/models` を呼び出し、通常の Gemini CLI モデルと `-search` バリアントの両方が存在することを確認します。
2. `gemini-xxx-search` をリクエストし、上流リクエストのモデル名が `gemini-xxx` であることを確認します。
3. 上流 Gemini CLI / Code Assist リクエスト本文に `googleSearch` が含まれることを確認します。
4. 通常の `gemini-xxx` をリクエストし、`googleSearch` が注入されないことを確認します。
5. デプロイ後、有効な Gemini CLI auth を使って実際の検索動作を検証します。
