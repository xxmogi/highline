# 仕様書: prompt / time / newline セグメント

## 概要

`highline` に 3 つの新機能を追加する。

| 機能 | セグメント名 | 概要 |
|------|-------------|------|
| プロンプト記号 | `prompt` | `$` / `#` / `%` をセグメントとして表示 |
| 時刻 | `time` | 現在時刻をセグメントとして表示 |
| 改行 | `newline` | プロンプトを複数行に分割する |

---

## ユースケース

### UC-1: プロンプト記号セグメント

**アクター**: シェルユーザー

**事前条件**: highline が `prompt` セグメントを含む設定で実行される

**基本フロー**:
1. `prompt` セグメントが `Render()` を呼ぶ
2. 現在のユーザーが root（UID=0）の場合 → `#` を表示
3. 非 root かつ shell=zsh の場合 → `%` を表示
4. 非 root かつ shell=bash（デフォルト）の場合 → `$` を表示
5. セグメントとして fg/bg 色付きで描画される

**設定**:
```json
{
  "segments": ["path", "git", "prompt"],
  "theme": {
    "prompt": {"fg": 15, "bg": 0}
  }
}
```

**デフォルト色**: FG=15（明るい白）、BG=0（黒）

---

### UC-2: 時刻セグメント

**アクター**: シェルユーザー

**事前条件**: highline が `time` セグメントを含む設定で実行される

**基本フロー**:
1. `time` セグメントが `Render()` を呼ぶ
2. `time.Now()` で現在時刻を取得し、フォーマット文字列で整形
3. セグメントとして fg/bg 色付きで描画される

**設定**:
```json
{
  "segments": ["time", "path"],
  "time": {"format": "15:04:05"}
}
```

**デフォルトフォーマット**: `"15:04"`（Go の time.Format 記法）
**デフォルト色**: FG=0（黒）、BG=7（白）

---

### UC-3: 改行セグメント

**アクター**: シェルユーザー

**事前条件**: highline が `newline` セグメントを含む設定で実行される

**基本フロー**:
1. `newline` セグメントを segments リストに含める
2. レンダラーが `newline` に到達した時点で、それまでの Powerline チェーンを閉じる
3. 改行文字 `\n` を出力する
4. 以降のセグメントは新しい行から描画を再開する

**例**:
```
segments: ["path", "git", "newline", "prompt"]
```
↓ 出力（イメージ）
```
 ~/work/highline  main ❯ \n$ ❯
```

**注意**: `newline` セグメントは色を持たず、Powerline セパレータも描画しない。
`newline` が最初・最後のセグメントの場合は空行扱い（特にエラーにはしない）。

---

## データ定義

### `PromptConfig`
```go
type PromptConfig struct {
    Shell string // "bash" or "zsh"
    FG    *int   // nil = 15
    BG    *int   // nil = 0
}
```

### `TimeConfig`
```go
type TimeConfig struct {
    Format string // Go time.Format string; "" = "15:04"
    FG     *int   // nil = 0
    BG     *int   // nil = 7
}
```

### `config.Config` への追加
```go
type Config struct {
    // 既存フィールド...
    Time TimeOptions `json:"time"`
}

type TimeOptions struct {
    Format string `json:"format"` // "" = "15:04"
}
```

---

## 制約・非機能要件

- `newline` セグメントはレンダラーが名前で識別し、`Render()` は呼ばれない
- `newline` をビルトインテーマに含める必要はない（色なし）
- `prompt` の記号はOSの実際のUID（`os.Getuid()`）で決定する
- `time` のフォーマット文字列は Go の `time.Format` 記法に従う（バリデーションなし）
- テストでは `time.Now()` の代わりに注入可能な clock を使う

---

## 受け入れ条件

1. `go test ./...` が全テスト通過
2. `go build ./cmd/highline` が成功
3. `./highline path newline prompt` で 2 行プロンプトが出力される
4. `./highline time path` で現在時刻が先頭セグメントに表示される
5. `./highline --config /dev/null prompt` で `$` が表示される（非rootの場合）
6. config.json に `"time": {"format": "2006-01-02"}` を書き、日付形式で表示される
