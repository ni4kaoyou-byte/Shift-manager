# 05_data-model.md — データモデル設計（MVP）

- プロダクト：飲食店向け シフト管理SaaS（個人開発MVP）
- 根拠ドキュメント：`docs/01_problem-users.md` / `docs/02_requirements.md` / `docs/04_ux-flows.md`
- 前提技術：`Supabase Postgres` + `Go API (gin + sqlc)`
- 更新：v0.1

---

## 0. 目的

- MVP機能（FR-01〜FR-20）を破綻なく支える最小データ構造を定義する
- テナント分離（店舗単位）と監査性（変更履歴）をDBレベルで担保する
- 将来拡張（多店舗/自動化）に耐える形で、今は過剰設計しない

---

## 1. モデリング方針

- IDは原則 `uuid`
- テナント境界は `store_id` で統一
- アプリユーザーは `auth.users`（Supabase）を正として、業務プロフィールは `app_users` で管理
- 公開後変更の証跡は `audit_logs` に必ず記録（FR-20）
- 削除は基本ソフトデリート（`disabled_at` など）で履歴保全

---

## 2. エンティティ一覧（MVP）

1. `app_users`  
用途：ユーザープロフィール（表示名など）

2. `stores`  
用途：店舗（MVPは1店舗運用だが、設計は多店舗対応可能にしておく）

3. `store_memberships`  
用途：店舗所属とロール管理（manager/staff）

4. `invites`  
用途：招待リンク/コード参加（FR-04）

5. `schedule_periods`  
用途：シフト期間（開始/終了、締切、下書き/公開）

6. `availability_submissions`  
用途：スタッフの希望提出ヘッダ（提出済み判定の基準）

7. `availability_entries`  
用途：希望提出の明細（日付・時間帯）

8. `shift_assignments`  
用途：店長が作成する実シフト（公開後更新も含む）

9. `coverage_targets`  
用途：必要人数テンプレ（FR-13: Should）

10. `change_requests`  
用途：欠勤/変更申請（FR-16, FR-17）

11. `notifications`  
用途：アプリ内通知（FR-19）

12. `audit_logs`  
用途：変更履歴（誰が/いつ/何を）（FR-20）

---

## 3. 関係（ERの骨子）

- `app_users` 1 - n `store_memberships`
- `stores` 1 - n `store_memberships`
- `stores` 1 - n `schedule_periods`
- `schedule_periods` 1 - n `availability_submissions`
- `availability_submissions` 1 - n `availability_entries`
- `schedule_periods` 1 - n `shift_assignments`
- `shift_assignments` 1 - n `change_requests`（原則）
- `stores` 1 - n `notifications`
- `stores` 1 - n `audit_logs`

---

## 4. テーブル定義（最小カラム）

## 4.1 ユーザー・店舗・権限

### `app_users`

- `id uuid pk`（=`auth.users.id`）
- `display_name text not null`
- `created_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`

### `stores`

- `id uuid pk`
- `name text not null`
- `created_by uuid not null fk -> app_users.id`
- `created_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`

### `store_memberships`

- `id uuid pk`
- `store_id uuid not null fk -> stores.id`
- `user_id uuid not null fk -> app_users.id`
- `role text not null check (role in ('manager','staff'))`
- `status text not null check (status in ('active','disabled')) default 'active'`
- `joined_at timestamptz not null default now()`
- `disabled_at timestamptz null`
- `unique(store_id, user_id)`

インデックス
- `idx_memberships_user (user_id)`
- `idx_memberships_store_role (store_id, role)`

### `invites`

- `id uuid pk`
- `store_id uuid not null fk -> stores.id`
- `code text not null unique`
- `expires_at timestamptz null`
- `max_uses int not null default 1`
- `used_count int not null default 0`
- `status text not null check (status in ('active','revoked','expired')) default 'active'`
- `created_by uuid not null fk -> app_users.id`
- `created_at timestamptz not null default now()`

---

## 4.2 シフト期間・希望・割当

### `schedule_periods`

- `id uuid pk`
- `store_id uuid not null fk -> stores.id`
- `name text not null`（例: 2026-02 第3週）
- `start_date date not null`
- `end_date date not null`
- `submission_deadline timestamptz not null`
- `status text not null check (status in ('draft','published')) default 'draft'`
- `published_at timestamptz null`
- `published_by uuid null fk -> app_users.id`
- `created_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`
- `check (start_date <= end_date)`

インデックス
- `idx_periods_store_dates (store_id, start_date desc)`
- `idx_periods_store_status (store_id, status)`

### `availability_submissions`

- `id uuid pk`
- `period_id uuid not null fk -> schedule_periods.id`
- `staff_user_id uuid not null fk -> app_users.id`
- `submitted_by uuid not null fk -> app_users.id`（本人 or 店長代理入力）
- `submitted_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`
- `unique(period_id, staff_user_id)`

インデックス
- `idx_avail_submission_period (period_id)`

### `availability_entries`

- `id uuid pk`
- `submission_id uuid not null fk -> availability_submissions.id on delete cascade`
- `work_date date not null`
- `start_time time not null`
- `end_time time not null`
- `availability_type text not null check (availability_type in ('available','unavailable')) default 'available'`
- `check (start_time < end_time)`

インデックス
- `idx_avail_entries_submission_date (submission_id, work_date)`

### `shift_assignments`

- `id uuid pk`
- `period_id uuid not null fk -> schedule_periods.id`
- `store_id uuid not null fk -> stores.id`
- `staff_user_id uuid not null fk -> app_users.id`
- `work_date date not null`
- `start_time time not null`
- `end_time time not null`
- `created_by uuid not null fk -> app_users.id`
- `updated_by uuid not null fk -> app_users.id`
- `created_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`
- `check (start_time < end_time)`

インデックス
- `idx_assignments_period_date (period_id, work_date)`
- `idx_assignments_staff_date (staff_user_id, work_date)`
- `idx_assignments_store_period (store_id, period_id)`

補足
- MVPでは重複勤務チェックはAPI側で実施（DB制約は将来強化）

### `coverage_targets`（Should）

- `id uuid pk`
- `period_id uuid not null fk -> schedule_periods.id`
- `day_of_week int not null check (day_of_week between 0 and 6)`（0=Sunday）
- `start_time time not null`
- `end_time time not null`
- `required_count int not null check (required_count >= 0)`
- `unique(period_id, day_of_week, start_time, end_time)`

---

## 4.3 申請・通知・監査

### `change_requests`

- `id uuid pk`
- `store_id uuid not null fk -> stores.id`
- `period_id uuid not null fk -> schedule_periods.id`
- `assignment_id uuid not null fk -> shift_assignments.id`
- `requester_user_id uuid not null fk -> app_users.id`
- `request_type text not null check (request_type in ('absence','change'))`
- `reason text null`
- `status text not null check (status in ('pending','approved','rejected')) default 'pending'`
- `reviewed_by uuid null fk -> app_users.id`
- `reviewed_at timestamptz null`
- `decision_note text null`
- `created_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`

インデックス
- `idx_change_requests_store_status (store_id, status)`
- `idx_change_requests_requester (requester_user_id, created_at desc)`

### `notifications`

- `id uuid pk`
- `store_id uuid not null fk -> stores.id`
- `user_id uuid not null fk -> app_users.id`
- `type text not null check (type in ('deadline_reminder','schedule_published','schedule_changed','request_decided'))`
- `title text not null`
- `body text not null`
- `payload jsonb not null default '{}'::jsonb`
- `read_at timestamptz null`
- `created_at timestamptz not null default now()`

インデックス
- `idx_notifications_user_created (user_id, created_at desc)`
- `idx_notifications_unread (user_id) where read_at is null`

### `audit_logs`

- `id uuid pk`
- `store_id uuid not null fk -> stores.id`
- `actor_user_id uuid not null fk -> app_users.id`
- `entity_type text not null`（例: schedule_period / assignment / change_request）
- `entity_id uuid not null`
- `action text not null`（create / update / publish / approve / reject など）
- `before_data jsonb null`
- `after_data jsonb null`
- `created_at timestamptz not null default now()`

インデックス
- `idx_audit_logs_store_created (store_id, created_at desc)`
- `idx_audit_logs_entity (entity_type, entity_id)`

---

## 5. 要件マッピング（FR）

- FR-01/02：`app_users`, `store_memberships`
- FR-03/04：`stores`, `invites`, `store_memberships`
- FR-06/07/08：`schedule_periods`
- FR-09/10/11：`availability_submissions`, `availability_entries`
- FR-12/13/14/15/18：`shift_assignments`, `coverage_targets`, `schedule_periods`
- FR-16/17：`change_requests`
- FR-19：`notifications`
- FR-20：`audit_logs`

---

## 6. 代表クエリ（実装で最初に必要）

1. 期間の未提出スタッフ一覧
- `store_memberships (staff)` から `availability_submissions(period_id)` を左結合して抽出

2. スタッフの確定シフト取得
- `schedule_periods.status='published'` の `shift_assignments` を `staff_user_id` で抽出

3. 店長の未対応申請件数
- `change_requests` を `store_id + status='pending'` で集計

4. 履歴一覧（時系列）
- `audit_logs` を `store_id` で `created_at desc` 取得

---

## 7. RLS/権限の最低ルール（Supabase前提）

- スタッフは自分の `availability_submissions` / `change_requests` / `notifications` のみ参照可能
- 店長は所属 `store_id` の全業務データにアクセス可能
- `audit_logs` は店長のみ参照可能
- `store_id` 不一致はAPI層でも必ず拒否（二重防御）

---

## 8. MVPでやらない設計

- 高度な労務制約（連勤・休憩・深夜）をDB制約で完全表現
- 代打自動マッチング用の候補最適化テーブル
- 多店舗横断分析専用の集計マート

MVPは「回収→作成→公開→変更→履歴」が破綻なく回る最小モデルを優先する。
