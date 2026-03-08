# 仕様書: コンフィグファイルサポート

## 概要

`~/.config/highline/config.json` を読み込み、テーマ（各セグメントの色）と表示内容（セグメントの種類・順序・オプション）をコントロールできるようにする。
CLIフラグは引き続き使用可能で、設定の優先順位は **CLIフラグ > コンフィグファイル > デフォルト値** とする。

---

## ユースケース

### UC-04: コンフィグファイルによる設定

- **アクター**: ユーザー
- **事前条件**: `~/.config/highline/config.json` が存在する
- **基本フロー**:
  1. `highline` 起動時にコンフィグファイルを読み込む
  2. `segments` フィールドが設定されていれば、表示するセグメントとその順序をコンフィグから決定する
  3. `theme` フィールドが設定されていれば、各セグメントの色をコンフィグから適用する
  4. セグメント固有のオプション（`path.max_depth` 等）を反映する
- **代替フロー**:
  - コンフィグファイルが存在しない → デフォルト値で動作する（エラーにしない）
  - コンフィグファイルのパースに失敗する → エラーメッセージを stderr に出力してデフォルト値を使用する
- **事後条件**: コンフィグファイルの設定が反映されたプロンプトが出力される

### UC-05: CLIフラグによるコンフィグ上書き

- **アクター**: ユーザー
- **基本フロー**:
  1. `--shell` フラグが指定された場合、コンフィグの `shell` より優先する
  2. セグメント名が引数で指定された場合、コンフィグの `segments` より優先する
- **事後条件**: CLIフラグの値が反映される

### UC-06: コンフィグファイルパスの変更

- **アクター**: ユーザー
- **基本フロー**:
  1. `--config /path/to/config.json` フラグで任意のパスを指定できる
- **事後条件**: 指定したファイルが読み込まれる

---

## コンフィグファイル仕様

### ファイルパス（優先順位）

1. `--config` フラグで指定したパス
2. `$XDG_CONFIG_HOME/highline/config.json`（`XDG_CONFIG_HOME` 未設定時は `~/.config/highline/config.json`）

### フォーマット: JSON

```json
{
  "shell": "bash",
  "segments": ["path", "git", "kube"],
  "theme": {
    "path": { "fg": 7, "bg": 4 },
    "git":  { "fg": 0, "bg": 2 },
    "kube": { "fg": 7, "bg": 6 }
  },
  "path": {
    "max_depth": 4
  }
}
```

### フィールド定義

| フィールド | 型 | デフォルト | 説明 |
|-----------|-----|----------|------|
| `shell` | string | `"bash"` | 対象シェル (`bash` / `zsh`) |
| `segments` | []string | `["path","git","kube"]` | 表示するセグメントと順序 |
| `theme.<name>.fg` | int | セグメント既定値 | 前景色（ANSI 基本16色: 0-15） |
| `theme.<name>.bg` | int | セグメント既定値 | 背景色（ANSI 基本16色: 0-15） |
| `path.max_depth` | int | `4` | パスの最大表示深さ（これを超えると中間を省略） |

---

## データ定義

```go
// Config はコンフィグファイルのルート構造体
type Config struct {
    Shell    string                `json:"shell"`
    Segments []string              `json:"segments"`
    Theme    map[string]ColorTheme `json:"theme"`
    Path     PathOptions           `json:"path"`
}

// ColorTheme はセグメントごとの色設定
type ColorTheme struct {
    FG *int `json:"fg"` // nil = デフォルトを使用
    BG *int `json:"bg"`
}

// PathOptions は path セグメント固有の設定
type PathOptions struct {
    MaxDepth int `json:"max_depth"` // 0 = デフォルト(4)
}
```

---

## 制約・非機能要件

- 外部ライブラリ不使用（`encoding/json` のみ使用）
- コンフィグファイルが存在しない場合は無音でスキップ（エラーなし）
- パースエラー時は stderr に警告を出力し、デフォルトで動作する
- 色値は 0〜15 の範囲に制限する（範囲外はデフォルト値を使用）

---

## 受け入れ条件

| # | 条件 |
|---|------|
| 9  | コンフィグファイルが存在しない場合、デフォルト動作する |
| 10 | `segments` でセグメントの種類と順序を変更できる |
| 11 | `theme` で各セグメントの色を変更できる |
| 12 | `path.max_depth` でパスの省略深さを変更できる |
| 13 | `--config` フラグで任意のコンフィグファイルを指定できる |
| 14 | `--shell` / セグメント引数はコンフィグより優先される |
| 15 | `make test` で全テストが通る |
