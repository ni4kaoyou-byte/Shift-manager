# 03_tech-stack.md — 技術スタック（Go + React / 無料枠優先）

- プロダクト：飲食店向け シフト管理SaaS（個人開発MVP）
- 方針：**Go + React** を採用し、MVPは**無料枠中心で運用開始**する
- 対象：1店舗（店長 + スタッフ）
- 更新：v0.1

---

## 0. スタック選定方針

- MVPで最優先するのは「速く作る」よりも「運用が破綻しないこと」
- 無料枠で始め、利用増加時に段階的に有料化できる構成にする
- アプリ責務はGo APIに集約し、フロントはUI/状態管理に集中する

---

## 1. 採用スタック（MVP固定）

### フロントエンド

- `React 19 + TypeScript`
- `Vite`（ビルド/開発サーバ）
- `React Router`（画面遷移）
- `TanStack Query`（サーバ状態）
- `React Hook Form + Zod`（入力/バリデーション）
- `Tailwind CSS`（UIスタイリング）
- `ESLint + Prettier`（静的解析/フォーマット）

### バックエンド

- `Go 1.23+`
- `gin`（HTTPルータ）
- `sqlc + pgx`（型安全SQLアクセス）
- `go-playground/validator`（追加バリデーション）
- `zerolog`（構造化ログ）

### データ・認証

- `Supabase Auth`（メール+パスワード認証）
- `Supabase Postgres`（メインDB）
- `Supabase Storage`（必要時のみ、添付用途）

### インフラ（無料枠優先）

- フロント配信：`Cloudflare Pages (Free)`
- Go APIホスティング：`Render Web Service (Free)`
- DB/Auth：`Supabase Free`

### テスト

- フロント：`Vitest + Testing Library`
- バックエンド：`go test`
- E2E：`Playwright`

---

## 2. システム構成（MVP）

1. ユーザーはReactアプリ（Cloudflare Pages）にアクセス  
2. ログインはSupabase Authで実施し、JWTを取得  
3. ReactはJWT付きでGo APIへリクエスト  
4. Go APIはJWT検証し、RBAC/テナント境界をチェック  
5. Go APIがSupabase Postgresに対して読み書き  
6. 公開・変更確定などのイベントはアプリ内通知テーブルへ保存

---

## 3. 無料枠前提の注意点（重要）

### Render Free

- 15分アイドルでスリープするため、最初の1リクエストが遅くなる
- 本番用途には制約が多いので、MVP検証用と割り切る

### Supabase Free

- DBサイズ上限（500MB/プロジェクト）を超えるとread-onlyになる
- プロジェクトは無料で2つまで（個人開発なら十分）

### Cloudflare Pages Free

- 静的配信は無料で扱いやすい
- Functionsを使う場合はWorkers無料枠の範囲で運用する

---

## 4. この構成が要件に合う理由

- FR-01/02（認証/ロール）  
  - Supabase Auth + Go側RBACで実装しやすい
- FR-03〜FR-20（業務ロジック）  
  - 複雑なシフトルールはGo APIに閉じ込め、フロントを軽く保てる
- NFR-02（アクセス制御）  
  - JWT + 店舗IDスコープをAPIで一元チェックできる
- NFR-08（モバイル最適）  
  - React/TailwindでスマホUIを先行で整備しやすい

---

## 5. 実装ルール（MVP時点）

- DBアクセスは`sqlc`生成コード経由のみ（生SQL散在を防ぐ）
- APIは`/api/v1`配下でバージョニング
- 監査ログ（変更履歴）は公開後変更で必ず記録（削除APIは作らない）
- 通知は初期はアプリ内通知を優先（メールは後追い）
- フロントコードは `ESLint + Prettier` を必須化し、PR前に整形/lintを通す

---

## 6. 将来の拡張方針（有料化の順番）

1. APIをRender有料プランへ（スリープ解消）
2. SupabaseをProへ（DB容量・可用性強化）
3. メール通知導入（Resend等）
4. バックグラウンドジョブ導入（催促通知/定期処理）

---

## 7. 代替案と不採用理由（メモ）

- Next.jsフルスタック  
  - 実装は速いが、今回は「Goで業務ロジックを持つ」方針に合わせて不採用
- Firebase中心構成  
  - SQLベースのシフト/監査ログ設計と相性が弱く不採用
- 自前認証（Goのみ）  
  - MVP初期の実装コストが高く、Supabase Auth利用を優先
