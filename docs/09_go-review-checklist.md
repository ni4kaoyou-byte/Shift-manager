# 09_go-review-checklist.md — Goレビュー項目（軽量クリーン前提）

- 対象: `apps/api`（Go）
- 前提: `handler -> usecase -> domain` の依存方向を守る
- 目的: レビュー時に「設計崩れ」と「運用事故」を早期検知する

---

## 0. 最重要: 境界・依存方向（まずここだけで8割）

- [ ] 依存方向が `handler -> usecase -> domain` になっている
- [ ] `domain` / `usecase` が `gin` / `sqlc` / `pgx` / `supabase` / HTTP に依存していない
- [ ] SQL実行が `infrastructure` に閉じている
- [ ] repository interface が usecase 側にある
- [ ] domain にDB都合の型（`sql.NullString`, `pgtype` など）を持ち込んでいない

レビューの最初に import 依存を確認し、崩れていたら最優先で修正する。

---

## 1. レイヤ別レビュー観点

### 1.1 handler（HTTP I/O）

- [ ] I/O変換 + バリデーション + usecase呼び出し + レスポンス整形に責務が限定されている
- [ ] バリデーションエラー形式が統一されている（400 + フィールド情報）
- [ ] DTO（request/response）が domain と分離されている
- [ ] HTTPステータス判断が適切（401/403/404/409/400/500）
- [ ] `c.Request.Context()` 由来の `context.Context` を usecase に渡している

地雷:
- handlerに業務分岐が入り始める
- domain構造体に `json` タグを付けてHTTP層都合を持ち込む

### 1.2 usecase（機能単位の手続き）

- [ ] 1ユースケース = 1機能の粒度で肥大化していない
- [ ] 認可の二重防御がある（middleware + usecase）
- [ ] トランザクション境界が明確（例: 状態更新 + 監査ログ）
- [ ] N+1 / 無駄なround tripがない
- [ ] DBエラーをそのまま返さず、ユースケース意味へ変換している
- [ ] 競合しやすい操作にロック/再試行方針がある

MVP最低限:
- Txが必要なところだけTx
- 競合しやすい操作（承認、公開など）だけロック方針を決める

### 1.3 domain（業務ルール）

- [ ] 状態遷移をdomainメソッドで表現している（不正遷移防止）
- [ ] 不変条件をコンストラクタ/ファクトリで保証している
- [ ] ID/値オブジェクトで取り違いを防いでいる
- [ ] 時刻注入（`now` 引数やClock）でテスト可能性を確保している
- [ ] domainがHTTPや文言都合を知らない

地雷:
- 状態遷移がusecaseに散らばる
- domainが `sql.NullTime` などインフラ都合を持つ

### 1.4 repository（永続化の抽象）

- [ ] interfaceがusecaseに必要な粒度である
- [ ] メソッド名がDB都合ではなく意図で命名されている
- [ ] 返却型がdomain中心で、infrastructure型を漏らさない
- [ ] `context.Context` を受け取る

### 1.5 infrastructure（sqlc + pgx / supabase）

- [ ] sqlc生成コードをアプリ層が直接importしていない
- [ ] DBエラーをドメイン/ユースケースエラーへ変換している
- [ ] クエリ（index、範囲検索、JOIN）が要件に対して妥当
- [ ] null/optionalマッピング方針が統一されている
- [ ] 外部I/Oにタイムアウト/リトライ/ログ方針がある

---

## 2. 横断レビュー観点（運用で死なないため）

### 2.1 認証・認可

- [ ] middlewareで `user_id` / `role` / `store_id` を確実に取り出す
- [ ] usecaseで対象リソースのstore境界を検証している
- [ ] role分岐が関数化/ポリシー化されている

### 2.2 エラーハンドリング

- [ ] レイヤを跨いでも意味が壊れないようにラップされている（`%w`）
- [ ] エラー→HTTPマッピングが集約されている
- [ ] 予期せぬエラーは安全なレスポンス + 十分な内部ログ

### 2.3 トランザクション

- [ ] Tx必須ケース（監査ログ、状態更新、関連更新）が明示されている
- [ ] usecaseがTx境界を制御できる（`WithTx` / `TxManager` 等）
- [ ] Tx内で外部APIを呼ばない

### 2.4 ログ・監査ログ

- [ ] `request_id` / `user_id` / `store_id` がログにある
- [ ] 監査ログの記録方針（成功時のみ/試行含む）が定義されている

### 2.5 テスト戦略

- [ ] domainユニットテスト（状態遷移/不変条件）
- [ ] usecaseテスト（repo差し替え）
- [ ] infra統合テスト（sqlcクエリ破損検知）

### 2.6 パフォーマンス最低限

- [ ] N+1の芽がない
- [ ] ページング/ソートの責務が明確

---

## 3. レビューの進め方（順番固定）

1. 依存方向/境界（import / package / interface位置）
2. 認可二重防御（store境界 + role）
3. Tx境界（承認・公開・監査ログ）
4. エラー設計（ドメイン/ユースケースエラー -> HTTP）
5. domain状態遷移（if散在を防ぐ）
6. 細部（命名、可読性、テスト）
