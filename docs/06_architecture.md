# 06_architecture.md — アーキテクチャ設計（MVP）

- プロダクト：飲食店向け シフト管理SaaS（個人開発MVP）
- 根拠ドキュメント：`docs/02_requirements.md` / `docs/03_tech-stack.md` / `docs/04_ux-flows.md` / `docs/05_data-model.md`
- 更新：v0.1

---

## 0. 結論（このプロジェクトの採用方針）

- バックエンド（Go）：**軽量クリーンアーキテクチャ**
- フロントエンド（React）：**feature-based architecture**
- リポジトリ構成：**monorepo（同一リポジトリ）**
- ディレクトリ分離：**`apps/api` と `apps/web` を分ける**
- Docker方針：**MVPはGo API中心にDocker化、Reactはローカル開発優先**

---

## 1. システム全体像

1. ユーザーは `apps/web`（React SPA）にアクセス
2. 認証は Supabase Auth で行い、JWTを取得
3. フロントはJWTを付与して `apps/api`（Go）へリクエスト
4. Go APIは認証・認可（role/store_id）を検証
5. Go APIがPostgres（Supabase）へアクセス
6. 公開/変更確定イベントは `notifications` と `audit_logs` に保存

---

## 2. リポジトリ・ディレクトリ構成

```txt
shift-manager/
  apps/
    api/                        # Go API
      cmd/server/               # エントリポイント
      internal/
        domain/                 # Entity / ValueObject / Domain Service
        usecase/                # アプリケーションサービス（ユースケース）
        handler/                # HTTPハンドラ（gin）
        repository/             # Repository interface
        infrastructure/         # DB実装、外部サービス実装
        middleware/             # 認証・認可・ロギング
        config/                 # 設定読込
    web/                        # React
      src/
        app/                    # ルータ、Provider、初期化
        features/               # 機能単位（auth/schedule/request...）
        shared/                 # 共通UI、API client、utils、types
  db/
    migrations/
    schema.sql
  docs/
```

---

## 3. Go APIアーキテクチャ（軽量クリーン）

### 3.1 依存方向

- `handler -> usecase -> domain`
- `infrastructure` は `repository interface` を実装
- `domain/usecase` は `gin` やSQL実装詳細に依存しない

### 3.2 レイヤ責務

- `domain`
  - 業務ルールの中心（期間状態、申請状態遷移など）
- `usecase`
  - 画面/機能単位の処理（例：シフト公開、欠勤申請承認）
- `handler`
  - HTTP I/O（JSON変換、バリデーション、レスポンス整形）
- `repository`
  - 永続化インターフェース
- `infrastructure`
  - `sqlc + pgx` の具体実装、Supabase連携

### 3.3 MVPで守る境界

- SQLは `infrastructure` のみで実行
- `usecase` から直接SQLを呼ばない
- 認可チェック（role/store_id）は `middleware + usecase` の二重防御

---

## 4. Reactアーキテクチャ（feature-based）

### 4.1 基本方針

- `features` 配下に機能を閉じ込める
- 画面は feature の組み合わせで構築
- API通信は feature 内 `api/`、表示は `components/`、状態は `hooks/` に整理

### 4.2 例

```txt
src/features/
  auth/
    api/
    components/
    hooks/
    pages/
  schedule/
    api/
    components/
    hooks/
    pages/
  change-request/
    api/
    components/
    hooks/
    pages/
```

### 4.3 共有領域

- `shared/ui`: 共通UI部品
- `shared/api`: HTTP client、token attach、error normalize
- `shared/types`: API DTO / 共通型

---

## 5. API・認証・認可

### 5.1 API設計

- ベースパス：`/api/v1`
- リソース中心（periods / availabilities / assignments / change-requests / notifications）
- 監査対象操作（公開、承認/却下、公開後変更）は必ず `audit_logs` を記録

### 5.2 認証

- Supabase JWTを `Authorization: Bearer` で受け取る
- APIで署名・期限・subを検証

### 5.3 認可

- `store_memberships` の `role` と `status` で制御
- 店長のみ：公開、割当編集、申請承認、履歴閲覧
- スタッフ：希望提出、自分の申請・通知・確定シフト閲覧

---

## 6. データアクセス方針

- DBスキーマの正本：`db/migrations/*.sql`
- `sqlc` 入力スキーマ：`db/schema.sql`
- クエリは `sqlc` 管理で型安全化
- トランザクションは usecase単位で明示実行（公開処理、申請承認処理）

---

## 7. Docker方針（MVP）

### 7.1 採用方針

- Go API：Docker化する（実行環境再現、デプロイ容易化）
- React：ローカル開発（`npm run dev`）を基本にする
- 全面Docker（web+api+db）はMVP初期では必須にしない

### 7.2 理由

- 開発速度を落としすぎず、API側の再現性だけ先に確保できる
- Reactホットリロード体験を維持しやすい
- DBは当面Supabaseを使うため、ローカルDB常駐を必須にしない

### 7.3 将来拡張

- CIやローカル再現性が必要になった時点で `docker-compose` を導入
- 必要なら `web` / `api` / `mock services` をcompose化

---

## 8. 非機能への対応

- NFR-02（アクセス制御）  
  - JWT検証 + store境界 + role境界の二重防御
- NFR-05（監査ログ保全）  
  - `audit_logs` を削除不可運用、公開後変更時は必須記録
- NFR-08（モバイル最適）  
  - Reactのスタッフ導線をスマホ優先で設計
- NFR-12（保守性）  
  - レイヤ責務を固定し、設計変更はADRで管理

---

## 9. 今はやらないこと（MVP外）

- イベント駆動の複雑な非同期基盤
- 厳密なCQRS分離
- マイクロサービス分割
- 複数フロント（管理画面別アプリ等）への分離

MVPでは「単一API + 単一Web」で運用の確実性を優先する。
