# Paste Formatted Date

> **これは日本語版（正本）です。**
> 英語版（参照）は [README.md](README.md) を参照してください。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/ci.yml)

Alfred 5 から今日の日付を複数のフォーマットで生成・貼り付けするワークフロー。

## 使い方

Alfred で `date` と入力するとフォーマット一覧が表示されます。選択すると自動的に貼り付けられます。

```
date             — フォーマット一覧を表示
date <filter>    — フォーマット名や値で絞り込み（例: "ISO", "YYYY", "unix"）
date config      — 設定の確認 / リセット
date help        — コマンド一覧を表示
```

### 利用可能なフォーマット

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

## 動作環境

- Alfred 5（Script Filter には Powerpack が必要）
- Python 3.9+

## インストール

[Releases](https://github.com/y-marui/alfred-paste-formatted-date/releases) から最新の `.alfredworkflow` をダウンロードしてダブルクリックでインストールします。

## 開発

```bash
# 開発用依存関係をインストール
make install

# Alfred をローカルでシミュレート
make run Q=""
make run Q="ISO"

# テストを実行
make test

# ワークフローパッケージをビルド
make build
# → dist/alfred-paste-formatted-date-0.1.0.alfredworkflow
```

## プロジェクト構成

```
alfred-paste-formatted-date/
├── src/
│   ├── alfred/         # Alfred SDK (response, router, cache, config, logger, safe_run)
│   └── app/            # アプリケーション層 (commands)
├── workflow/           # Alfred パッケージ (info.plist, scripts/entry.py, vendor/)
├── tests/              # pytest テストスイート
├── scripts/            # build.sh, dev.sh, release.sh, vendor.sh
└── docs/               # アーキテクチャ・開発ドキュメント
```

## サポート

このワークフローが役に立ったら、サポートしていただけると嬉しいです。

- [Buy Me a Coffee](https://www.buymeacoffee.com/y.marui)
- [GitHub Sponsors](https://github.com/sponsors/y-marui)

## ライセンス

MIT — [LICENSE](LICENSE) を参照

---

*この文書には英語版（参照版）[README.md](README.md) があります。編集時は同一コミットで更新してください。*
