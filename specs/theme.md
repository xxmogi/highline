# 仕様書: ビルトインテーマサポート

## 概要

Monokai・Solarized など一般的なエディタテーマをワンステップで適用できる組み込みテーマ機能を追加する。
`--theme` フラグまたはコンフィグの `theme_name` フィールドでテーマ名を指定するだけで、全セグメントの色が切り替わる。
既存の `theme` フィールドによる個別上書きは引き続き有効で、テーマカラーより優先される。

---

## ユースケース

### UC-07: ビルトインテーマの適用

- **アクター**: ユーザー
- **事前条件**: `highline` バイナリが PATH に存在する
- **基本フロー**:
  1. ユーザーが `--theme monokai` のように CLI フラグでテーマ名を指定する
  2. 指定テーマの色定義が全セグメントに適用される
  3. Powerline スタイルのプロンプトが出力される
- **代替フロー**:
  - 存在しないテーマ名を指定 → stderr に警告を出力し、デフォルト色で動作する
- **事後条件**: 指定テーマが反映されたプロンプトが出力される

### UC-08: コンフィグによるテーマ指定

- **アクター**: ユーザー
- **基本フロー**:
  1. `config.json` に `"theme_name": "solarized-dark"` を記述する
  2. `highline` 起動時にテーマが自動で適用される
- **事後条件**: コンフィグのテーマが反映される

### UC-09: テーマと個別色設定の併用

- **アクター**: ユーザー
- **基本フロー**:
  1. `theme_name` でベーステーマを指定する
  2. `theme.<segment>` で特定セグメントの色だけ上書きする
- **事後条件**: テーマ色をベースに、個別設定が上書きされた状態でプロンプトが表示される

---

## テーマ一覧

| テーマ名         | 説明 |
|-----------------|------|
| `monokai`       | Sublime Text の Monokai テーマ。鮮やかな緑・黄・マゼンタ |
| `solarized-dark`  | Solarized Dark。落ち着いたブルーグリーン系 |
| `solarized-light` | Solarized Light。明るい背景向け |
| `nord`          | Arctic Ice Studio の Nord テーマ。青みがかったグレー系 |
| `dracula`       | Dracula テーマ。ダークパープル・グリーン系 |

### 各テーマのカラー定義（ANSI 基本16色）

| テーマ              | セグメント | FG          | BG           |
|--------------------|-----------|-------------|--------------|
| `monokai`          | path      | 0 (Black)   | 10 (BrightGreen)   |
|                    | git       | 0 (Black)   | 11 (BrightYellow)  |
|                    | kube      | 0 (Black)   | 5  (Magenta)       |
| `solarized-dark`   | path      | 15 (BrightWhite) | 4 (Blue)    |
|                    | git       | 8 (BrightBlack)  | 2 (Green)   |
|                    | kube      | 8 (BrightBlack)  | 6 (Cyan)    |
| `solarized-light`  | path      | 8 (BrightBlack)  | 3 (Yellow)  |
|                    | git       | 8 (BrightBlack)  | 2 (Green)   |
|                    | kube      | 8 (BrightBlack)  | 6 (Cyan)    |
| `nord`             | path      | 15 (BrightWhite) | 4 (Blue)    |
|                    | git       | 15 (BrightWhite) | 12 (BrightBlue)  |
|                    | kube      | 8 (BrightBlack)  | 6 (Cyan)    |
| `dracula`          | path      | 15 (BrightWhite) | 5 (Magenta) |
|                    | git       | 0 (Black)        | 10 (BrightGreen) |
|                    | kube      | 15 (BrightWhite) | 12 (BrightBlue)  |

---

## 設定例

### CLI フラグ

```bash
highline --theme monokai
highline --theme solarized-dark --shell zsh
```

### config.json

```json
{
  "theme_name": "nord",
  "theme": {
    "path": { "bg": 12 }
  }
}
```

`theme` による個別設定は `theme_name` による色より優先される。

---

## データ定義

```go
// internal/theme/theme.go

// Palette はセグメント1つ分の色ペア。
type Palette struct {
    FG int
    BG int
}

// Theme はセグメント名 → Palette のマップ。
type Theme map[string]Palette

// Get はテーマ名（大文字小文字不問）からビルトインテーマを返す。
// 見つからない場合は false を返す。
func Get(name string) (Theme, bool)

// Names はビルトインテーマ名の一覧を返す（ソート済み）。
func Names() []string
```

```go
// internal/config/config.go への追加フィールド
type Config struct {
    Shell     string                `json:"shell"`
    ThemeName string                `json:"theme_name"` // 追加
    Segments  []string              `json:"segments"`
    Theme     map[string]ColorTheme `json:"theme"`
    Path      PathOptions           `json:"path"`
}
```

---

## 色解決の優先順位

```
CLI --theme > config.theme_name > config.theme（個別設定） > セグメントデフォルト
```

ただし `config.theme` の個別上書きは名前付きテーマより優先されるため、実際の解決順は：

1. セグメントデフォルト色を初期値とする
2. 名前付きテーマの色で上書きする（`--theme` > `config.theme_name`）
3. `config.theme.<segment>` の個別指定でさらに上書きする

---

## 制約・非機能要件

- 外部ライブラリ不使用
- 色値は ANSI 基本16色（0〜15）に限定する
- 存在しないテーマ名はエラーにせず、stderr に警告を出力してデフォルト動作する
- `internal/theme` パッケージは `internal/config` や `internal/segment` に依存しない

---

## 受け入れ条件

| # | 条件 |
|---|------|
| 16 | `--theme monokai` でMonokaiテーマが適用される |
| 17 | `config.json` の `theme_name` でテーマが適用される |
| 18 | `--theme` は `config.theme_name` より優先される |
| 19 | `config.theme` の個別設定はテーマ色より優先される |
| 20 | 存在しないテーマ名を指定しても stderr 警告のみでクラッシュしない |
| 21 | `make test` で全テストが通る |
