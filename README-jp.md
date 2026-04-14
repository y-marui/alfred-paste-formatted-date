# Alfred Workflow Template

> **これは日本語版（正本）です。**
> 英語版（参照）は [README.md](README.md) を参照してください。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
<!-- TODO: CI バッジのリンク先と画像 URL（href と src の両方）を自分のリポジトリ URL に書き換えること。「このテンプレートのカスタマイズ手順」を参照。 -->
[![CI](https://github.com/y-marui/alfred-workflow-template/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-workflow-template/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-workflow-template/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-workflow-template/actions/workflows/dev-charter-check.yml)

| 項目 | 内容 |
|---|---|
| 開発対象 | Alfred 5 Script Filter ワークフロー |
| 開発環境 | 個人〜小規模チーム（1〜3人） |
| 主言語 | 英語（OSS） |
| ライセンス | MIT |
| 動作環境 | Python 3.9+, Alfred 5 |
| AI ツール | Claude Code / GitHub Copilot / Gemini CLI |

> Alfred 5 Script Filter ワークフローのプロダクションレディなテンプレート。
> 10分で開発を開始できます。

## Features

- ✅ **レイヤードアーキテクチャ** — Alfred 境界とビジネスロジックを分離
- ✅ **軽量 Alfred SDK** — レスポンスビルダー、ルーター、キャッシュ、設定、ロガー
- ✅ **コマンドベース UX** — `wf search`、`wf open`、`wf config`、`wf help`
- ✅ **フルテストスイート** — pytest で Alfred なしにテスト実行可能
- ✅ **CI/CD** — GitHub Actions でリント・テスト・ビルド・リリースを自動化
- ✅ **ベンダーパッケージング** — サードパーティ依存を `vendor/` にバンドル
- ✅ **AI 対応** — `AI_CONTEXT.md` + `CLAUDE.md` で AI アシスタントのコンテキストを管理

## Requirements

- Alfred 5（Script Filter には Powerpack が必要）
- Python 3.9+
- [pre-commit](https://pre-commit.com/)（セキュリティフック用）

## Quick Start (developers)

```bash
git clone https://github.com/yourname/alfred-workflow-template
cd alfred-workflow-template

# 開発用依存関係をインストール
make install

# Alfred をローカルでシミュレート
make run Q="search foo"
make run Q="help"

# テストを実行
make test

# ワークフローパッケージをビルド
make build
# → dist/workflow-template-0.1.0.alfredworkflow
```

`dist/*.alfredworkflow` をダブルクリックして Alfred にインストールします。

## Usage

```
wf <query>           検索（デフォルト）
wf search <query>    検索
wf open <name>       ショートカットを開く
wf config            設定の確認 / リセット
wf help              コマンド一覧を表示
```

## Project Structure

```
alfred-workflow-template/
├── src/
│   ├── alfred/         # Alfred SDK (response, router, cache, config, logger, safe_run)
│   └── app/            # アプリケーション層 (commands, services, clients)
├── workflow/           # Alfred パッケージ (info.plist, scripts/entry.py, vendor/)
├── tests/              # pytest テストスイート
├── scripts/            # build.sh, dev.sh, release.sh, vendor.sh
└── docs/               # アーキテクチャ・開発・利用ドキュメント
```

## Documentation

| ドキュメント | 内容 |
|---|---|
| [docs/architecture.md](docs/architecture.md) | アーキテクチャ全体設計 |
| [docs/development.md](docs/development.md) | コマンド追加・依存関係管理・リリース手順 |
| [docs/usage.md](docs/usage.md) | エンドユーザー向け利用ガイド |

## AI-Assisted Development

このテンプレートは AI 支援開発に対応しています。

| ツール | 役割 |
|---|---|
| Claude Code | アーキテクチャ設計・大規模変更・リファクタリング |
| GitHub Copilot | バグ修正・細かな実装・単体テスト作成 |
| Gemini CLI | ドキュメント管理 |

セッションコンテキスト: [`AI_CONTEXT.md`](AI_CONTEXT.md)、[`CLAUDE.md`](CLAUDE.md)

## Customizing This Template

1. `workflow/info.plist` を編集:
   - `bundleid` を自分のバンドル ID に変更（例: `com.yourname.workflowname`）
   - キーワード（`wf`）を自分のトリガーキーワードに変更
   - `uuidgen` で生成した UUID に置き換え
2. `src/app/clients/api_client.py` を実際の API クライアントに置き換え
3. `pyproject.toml` のワークフロー名を更新
4. `src/app/commands/open_cmd.py` のショートカットを更新
5. `workflow/icon.png` を追加

## Release

```bash
# 1. pyproject.toml のバージョンを更新
# 2. タグを付けてプッシュ
git tag v1.2.3
git push --tags
# GitHub Actions が .alfredworkflow をビルドして GitHub Release を作成
```

## Support

このテンプレートが役に立ったら、サポートしていただけると嬉しいです。

- [Buy Me a Coffee](https://www.buymeacoffee.com/YOUR_USERNAME)
- [GitHub Sponsors](https://github.com/sponsors/YOUR_USERNAME)

## License

MIT — [LICENSE](LICENSE) を参照

---

*この文書には英語版（参照版）[README.md](README.md) があります。編集時は同一コミットで更新してください。*
