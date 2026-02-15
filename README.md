# Shift Manager

飲食店向けシフト管理SaaS（個人開発MVP）の設計・実装リポジトリです。  
MVPの目的は「希望回収 → 作成 → 公開 → 変更対応」を破綻なく回すことです。

## 現在の状態

- ドキュメント主導で仕様/設計を整備中
- 実装方針
  - API先行（Go）
  - 画面は実API接続で並行実装（React）
  - 縦スライス順で完成させる

## 開発バージョン方針（BL-002）

- Go: `1.23+`（現状は `1.25.x` で確認）
- Node.js: `20+`（現状は `23.x` で確認）
- npm: `10+`

## ドキュメント

- 入口: `docs/00_index.md`
- 要件: `docs/02_requirements.md`
- 技術スタック: `docs/03_tech-stack.md`
- アーキテクチャ: `docs/06_architecture.md`
- 実装バックログ: `docs/07_backlog.md`
- 開発規約: `docs/08_dev-conventions.md`
- テスト戦略: `docs/10_testing.md`

## リポジトリ構成（BL-001時点）

```txt
shift-manager/
  apps/
    api/                        # Go API
      cmd/
        server/
      internal/
        domain/
        usecase/
        handler/
        repository/
        infrastructure/
        middleware/
        config/
    web/                        # React
      src/
        app/
        features/
        shared/
  db/
    migrations/
    schema.sql
  docs/
```

## 環境構築（現時点）

このリポジトリは現在、実装準備段階です。  
先にドキュメント確認とDB定義確認を行ってください。

1. リポジトリを取得
```bash
git clone https://github.com/ni4kaoyou-byte/Shift-manager.git
cd Shift-manager
```

2. ドキュメントを確認
```bash
cat docs/00_index.md
```

3. DB定義を確認
```bash
ls db/migrations
cat db/schema.sql
```

## 環境構築（実装フェーズ）

以下は実装フェーズで運用する手順です。

1. 必須ツールを準備
- Go 1.23+
- Node.js 20+
- npm 10+
- PostgreSQLクライアント（必要に応じて）

2. 環境変数を作成
```bash
cp apps/api/.env.example apps/api/.env
cp apps/web/.env.example apps/web/.env
```

APIの必須環境変数（未設定だと起動時にエラー）
- `DATABASE_URL`
- `SUPABASE_URL`
- `SUPABASE_JWT_SECRET`

3. DB migration適用
```bash
cd apps/api
export DATABASE_URL='postgres://postgres:postgres@localhost:5432/shift_manager?sslmode=disable'
make db-up
```

4. sqlcコード生成
```bash
cd apps/api
make sqlc-generate
```

5. API起動
```bash
cd apps/api
make run
```

6. Web起動
```bash
cd apps/web
npm install
npm run dev
```

## DB運用ルール（BL-004）

- DBスキーマ変更は `db/migrations/*.sql` を正本として追加する
- migrationを追加/変更したら `db/schema.sql` も必ず同期する
- 同期後に `cd apps/api && make sqlc-generate` を実行する
- PR前チェック
  - migration適用手順がREADMEと一致している
  - `make sqlc-generate` が成功する

## よく使うコマンド

### API（Go）
```bash
cd apps/api
make run
make test
make lint
make fmt
```

### Web（React）
```bash
cd apps/web
npm run dev
npm run lint
npm run build
npm run test
```

## 開発ルール（抜粋）

- `main` 直push禁止、必ずブランチ + PR
- コミットメッセージは日本語
- フロントは `ESLint + Prettier` 必須
- 仕様変更時は docs/ADR を先に更新

## ライセンス

未定
