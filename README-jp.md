# Paste Formatted Date

> **これは日本語版（正本）です。**
> 英語版（参照）は [README.md](README.md) を参照してください。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

Alfred 5 から今日の日付を複数のフォーマットで生成・貼り付けするワークフロー。

## Usage

Alfred で `date` と入力するとフォーマット一覧が表示されます。選択すると自動的に貼り付けられます。

```
date             — フォーマット一覧を表示
date <filter>    — フォーマット名や値で絞り込み（例: "ISO", "YYYY", "unix"）
date config      — 設定の確認 / リセット
date help        — コマンド一覧を表示
```

### Available formats

| フォーマット | 例 |
|---|---|
| YYYYMMDD | 20260414 |
| YYMMDD | 260414 |
| YYYY-MM-DD | 2026-04-14 |
| YYYY/MM/DD | 2026/04/14 |
| MM/DD/YYYY | 04/14/2026 |
| DD/MM/YYYY | 14/04/2026 |
| MMM DD, YYYY | Apr 14, 2026 |
| MMMM DD, YYYY | April 14, 2026 |
| YYYY-MM-DDThh:mm:ss | 2026-04-14T12:00:00 |
| Unix timestamp | 1744588800 |

## Requirements

- Alfred 5（Script Filter には Powerpack が必要）

## Installation

[Releases](https://github.com/y-marui/alfred-paste-formatted-date/releases) から最新の `.alfredworkflow` をダウンロードしてダブルクリックでインストールします。

## Development

Go が必要です（ツールチェーンのバージョンは `go.mod` を参照）。開発フロー全体は [DEVELOPING.md](DEVELOPING.md) を参照してください。

```bash
# Alfred をローカルでシミュレート
go run ./cmd/paste-formatted-date-alfred ""
go run ./cmd/paste-formatted-date-alfred "ISO"

# テストを実行
make test

# ワークフローパッケージをビルド
make build-workflow
# → dist/paste-formatted-date-0.1.0.alfredworkflow
```

## Project Structure

```
alfred-paste-formatted-date/
├── cmd/
│   └── paste-formatted-date-alfred/  # Alfred が実行するバイナリ
├── internal/
│   ├── datecmd/         # コマンドディスパッチ (date / config / help)
│   ├── dateresolve/     # クエリ → 対象日付の解決
│   ├── datefmt/         # 日付フォーマット一覧・整形
│   ├── wfconfig/        # 永続設定ストア
│   └── scriptfilter/    # Alfred Script Filter JSON 型
├── workflow/            # Alfred パッケージ (info.plist, icon.png)
└── docs/                # アーキテクチャドキュメント
```

## License

MIT — [LICENSE](LICENSE) を参照

---

*この文書には英語版（参照版）[README.md](README.md) があります。編集時は同一コミットで更新してください。*
