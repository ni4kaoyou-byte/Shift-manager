# 00_index.md — ドキュメント索引

- プロダクト：飲食店向け シフト管理SaaS（個人開発MVP）
- 目的：ドキュメントの入口を一本化し、参照順と責務を明確にする
- 更新：v0.1

---

## 1. 読む順番（推奨）

1. `docs/01_problem-users.md`  
   何の課題を、誰のために解くか（問題定義/JTBD）

2. `docs/02_requirements.md`  
   MVPでやること・やらないこと（FR/NFR）

3. `docs/03_tech-stack.md`  
   技術選定（Go + React + Supabase + 無料枠方針）

4. `docs/04_ux-flows.md`  
   必要画面と主要ユーザーフロー

5. `docs/05_data-model.md`  
   エンティティ、テーブル設計、要件マッピング

6. `docs/06_architecture.md`  
   API/Webの構成方針、レイヤ責務、Docker方針

7. `docs/07_backlog.md`  
   API先行 + 縦スライスの実装順序

8. `docs/08_dev-conventions.md`  
   Git/PR、実装規約、品質規約、DoD

9. `docs/10_testing.md`  
   テスト戦略、E2E必須シナリオ、リリース前チェック

---

## 2. ADR一覧

- `docs/adr/0001-auth.md`  
  認証/認可の設計判断（作成予定）

- `docs/adr/0003-data-model-strategy.md`  
  データモデル戦略（store_id境界、監査ログ必須、migration正本）

---

## 3. DB関連ドキュメント

- `db/migrations/0001_init.up.sql` / `db/migrations/0001_init.down.sql`  
  実行用DDL（正本）

- `db/schema.sql`  
  sqlc用の集約スキーマ

---

## 4. 参照ルール

- 設計・方針を変更した場合
  1. 先に該当ドキュメントを更新
  2. 必要ならADRを追加
  3. 実装と整合するよう `db/schema.sql` / migration を同期

- 実装中に迷った場合の優先順位
  1. `docs/02_requirements.md`
  2. `docs/04_ux-flows.md`
  3. `docs/06_architecture.md`
  4. `docs/08_dev-conventions.md`
