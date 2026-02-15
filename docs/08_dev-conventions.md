# 08_dev-conventions.md — 開発規約（MVP）

- プロダクト：飲食店向け シフト管理SaaS（個人開発MVP）
- 対象：`apps/api`（Go） / `apps/web`（React） / `db`
- 更新：v0.1

---

## 0. 目的

- 実装時の判断を統一し、手戻りを減らす
- API先行 + 縦スライス開発を崩さない
- main直pushを防ぎ、PRベース運用に固定する

---

## 1. Git運用ルール

- `main` 直pushは禁止（例外なし）
- 作業は必ずブランチで行う
- ブランチ命名
  - `feature/<topic>`
  - `fix/<topic>`
  - `chore/<topic>`
- コミットメッセージは日本語で「何をしたか」を明確に書く
  - 例: `認証ミドルウェアを追加`
  - 例: `希望提出APIのバリデーションを修正`
- PRでマージするまで完了扱いにしない

---

## 2. 実装順序ルール（固定）

- API先行で薄い骨組みを作る
  - `auth`, `membership`, `period`, `availability`, `assignment`, `change_request`, `audit_log`
- 画面は実APIに接続して並行実装する（モック禁止）
- 縦スライスの順序を守る
  1. ログイン
  2. 参加
  3. 希望提出
  4. 公開
  5. 変更申請
  6. 履歴

---

## 3. アーキテクチャ規約

### Go（軽量クリーンアーキテクチャ）

- 依存方向は `handler -> usecase -> domain`
- SQL実行は `infrastructure` のみ
- `usecase` から `gin` / SQL実装詳細に依存しない
- `context.Context` はリクエスト境界で必須

### React（feature-based）

- `src/features/<feature>` 配下に機能を閉じ込める
- 共通化の基準
  - 2機能以上で使うもののみ `shared` へ移す
- API呼び出しは `features/*/api` or `shared/api` から行う

---

## 4. API規約

- ベースパスは `/api/v1`
- JSONは `camelCase` で統一
- 認可エラーの使い分け
  - `401`: 未認証（トークンなし/無効）
  - `403`: 認証済みだが権限不足
- バリデーションエラーは `400`
- 共通エラーレスポンス形式

```json
{
  "error": {
    "code": "string_code",
    "message": "human readable message"
  }
}
```

---

## 5. DB規約

- DB変更は必ず `db/migrations/*.sql` で行う
- `db/schema.sql` は sqlc 用集約スキーマとして同期する
- migration変更時の手順
  1. migration追加/更新
  2. `db/schema.sql` 同期
  3. sqlc再生成
  4. 影響APIのテスト更新
- テナント境界（`store_id`）チェックはAPIで必須

---

## 6. 品質規約（Lint/Format）

- フロントは `ESLint + Prettier` 必須
- Goは `gofmt`/`go vet` を必須
- PR前チェック
  - lint通過
  - test通過
  - 主要フローの手動確認

---

## 7. ログ・監査規約

- APIログに最低限含める項目
  - `request_id`
  - `user_id`（取得できる場合）
  - `store_id`（対象がある場合）
  - `path`, `status`, `latency_ms`
- `audit_logs` 記録対象
  - 期間公開
  - 公開後の割当変更
  - 申請承認/却下

---

## 8. 環境構築規約

- セットアップ手順はルート `README.md` を正本にする
- `.env.example` を更新せずに新規環境変数を追加しない
- 新メンバーがREADMEのみで起動できる状態を維持する

---

## 9. Definition of Done

- コード実装完了
- 受け入れ条件を満たす
- テスト（自動/手動）確認済み
- ドキュメント更新済み（必要時）
- PR作成済み、レビュー可能状態
