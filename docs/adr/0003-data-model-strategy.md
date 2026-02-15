# ADR 0003: Data Model Strategy（MVP）

- Date: 2026-02-15
- Status: Accepted
- Owners: shift-manager

---

## Context

- MVPの要件は `docs/02_requirements.md` の FR-01〜FR-20 を満たす必要がある。
- スコープは1店舗だが、将来の多店舗化を阻害しない設計が必要。
- 「公開後の変更履歴を残すこと」がプロダクト価値の中核（FR-20）。
- 技術スタックは Go API + Supabase Postgres + sqlc（`docs/03_tech-stack.md`）。

---

## Decision

1. テナント分離キーとして全業務データに `store_id` を持たせる  
2. 認証の正は Supabase `auth.users` とし、業務プロフィールは `public.app_users` で持つ  
3. 所属/権限は `store_memberships`（`manager` / `staff`）で表現する  
4. シフト業務は以下の分割で管理する  
   - 期間: `schedule_periods`  
   - 希望: `availability_submissions`, `availability_entries`  
   - 割当: `shift_assignments`  
   - 変更申請: `change_requests`  
5. 公開後の変更証跡は `audit_logs` に必ず保存し、削除機能は提供しない  
6. 通知はMVPではアプリ内通知 `notifications` を基準にする  
7. DDLの実行正本は `db/migrations/*.sql`、`db/schema.sql` は sqlc 用集約スキーマとする

---

## Consequences

### Positive

- 店舗単位のアクセス制御がAPI/RLS両方で実装しやすい
- 要件で重い「変更対応」「履歴保全」がDB構造で表現できる
- sqlcで型安全なクエリ生成がしやすい
- 1店舗MVPから多店舗化へ段階的に移行しやすい

### Negative

- `store_id` の整合性チェックが多く、実装時の注意点が増える
- 履歴記録（`audit_logs`）の運用コストが増える
- Supabase Auth と `app_users` の二重管理ポイントが発生する

---

## Alternatives Considered

1. 店舗概念を持たず単一テナント前提にする  
却下理由: MVP後の多店舗化で全面的な再設計が必要になるため。

2. 監査ログを持たず `updated_at` のみで運用する  
却下理由: FR-20（誰が/いつ/何を）を満たせない。

3. Firebase/NoSQL中心で設計する  
却下理由: シフト集計・履歴参照・SQL運用との相性が弱い。

---

## Follow-ups

1. `db/migrations/0001_init.up.sql` / `down.sql` を正として維持する  
2. `db/schema.sql` をマイグレーション変更時に同期する  
3. 次ADRで RLS方針（スタッフ/店長の境界）を明文化する  
4. 次ADRで監査ログ記録対象イベントを固定する

