-- sqlc aggregated schema
-- Source of truth for deployment is: db/migrations/*.sql
-- Keep this file in sync with the latest migration.

CREATE SCHEMA IF NOT EXISTS auth;

-- Supabase has auth.users, but sqlc parsing needs the relation locally.
CREATE TABLE IF NOT EXISTS auth.users (
  id uuid PRIMARY KEY
);

CREATE TABLE public.app_users (
  id uuid PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  display_name text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE public.stores (
  id uuid PRIMARY KEY,
  name text NOT NULL,
  created_by uuid NOT NULL REFERENCES public.app_users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE public.store_memberships (
  id uuid PRIMARY KEY,
  store_id uuid NOT NULL REFERENCES public.stores(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES public.app_users(id) ON DELETE CASCADE,
  role text NOT NULL CHECK (role IN ('manager', 'staff')),
  status text NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'disabled')),
  joined_at timestamptz NOT NULL DEFAULT now(),
  disabled_at timestamptz NULL,
  UNIQUE (store_id, user_id)
);

CREATE TABLE public.invites (
  id uuid PRIMARY KEY,
  store_id uuid NOT NULL REFERENCES public.stores(id) ON DELETE CASCADE,
  code text NOT NULL UNIQUE,
  expires_at timestamptz NULL,
  max_uses int NOT NULL DEFAULT 1 CHECK (max_uses > 0),
  used_count int NOT NULL DEFAULT 0 CHECK (used_count >= 0 AND used_count <= max_uses),
  status text NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'revoked', 'expired')),
  created_by uuid NOT NULL REFERENCES public.app_users(id),
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE public.schedule_periods (
  id uuid PRIMARY KEY,
  store_id uuid NOT NULL REFERENCES public.stores(id) ON DELETE CASCADE,
  name text NOT NULL,
  start_date date NOT NULL,
  end_date date NOT NULL,
  submission_deadline timestamptz NOT NULL,
  status text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published')),
  published_at timestamptz NULL,
  published_by uuid NULL REFERENCES public.app_users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (start_date <= end_date)
);

CREATE TABLE public.availability_submissions (
  id uuid PRIMARY KEY,
  period_id uuid NOT NULL REFERENCES public.schedule_periods(id) ON DELETE CASCADE,
  staff_user_id uuid NOT NULL REFERENCES public.app_users(id),
  submitted_by uuid NOT NULL REFERENCES public.app_users(id),
  submitted_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (period_id, staff_user_id)
);

CREATE TABLE public.availability_entries (
  id uuid PRIMARY KEY,
  submission_id uuid NOT NULL REFERENCES public.availability_submissions(id) ON DELETE CASCADE,
  work_date date NOT NULL,
  start_time time NOT NULL,
  end_time time NOT NULL,
  availability_type text NOT NULL DEFAULT 'available' CHECK (availability_type IN ('available', 'unavailable')),
  CHECK (start_time < end_time)
);

CREATE TABLE public.shift_assignments (
  id uuid PRIMARY KEY,
  period_id uuid NOT NULL REFERENCES public.schedule_periods(id) ON DELETE CASCADE,
  store_id uuid NOT NULL REFERENCES public.stores(id) ON DELETE CASCADE,
  staff_user_id uuid NOT NULL REFERENCES public.app_users(id),
  work_date date NOT NULL,
  start_time time NOT NULL,
  end_time time NOT NULL,
  created_by uuid NOT NULL REFERENCES public.app_users(id),
  updated_by uuid NOT NULL REFERENCES public.app_users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CHECK (start_time < end_time)
);

CREATE TABLE public.coverage_targets (
  id uuid PRIMARY KEY,
  period_id uuid NOT NULL REFERENCES public.schedule_periods(id) ON DELETE CASCADE,
  day_of_week int NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
  start_time time NOT NULL,
  end_time time NOT NULL,
  required_count int NOT NULL CHECK (required_count >= 0),
  CHECK (start_time < end_time),
  UNIQUE (period_id, day_of_week, start_time, end_time)
);

CREATE TABLE public.change_requests (
  id uuid PRIMARY KEY,
  store_id uuid NOT NULL REFERENCES public.stores(id) ON DELETE CASCADE,
  period_id uuid NOT NULL REFERENCES public.schedule_periods(id) ON DELETE CASCADE,
  assignment_id uuid NOT NULL REFERENCES public.shift_assignments(id),
  requester_user_id uuid NOT NULL REFERENCES public.app_users(id),
  request_type text NOT NULL CHECK (request_type IN ('absence', 'change')),
  reason text NULL,
  status text NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
  reviewed_by uuid NULL REFERENCES public.app_users(id),
  reviewed_at timestamptz NULL,
  decision_note text NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE public.notifications (
  id uuid PRIMARY KEY,
  store_id uuid NOT NULL REFERENCES public.stores(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES public.app_users(id) ON DELETE CASCADE,
  type text NOT NULL CHECK (type IN (
    'deadline_reminder',
    'schedule_published',
    'schedule_changed',
    'request_decided'
  )),
  title text NOT NULL,
  body text NOT NULL,
  payload jsonb NOT NULL DEFAULT '{}'::jsonb,
  read_at timestamptz NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE public.audit_logs (
  id uuid PRIMARY KEY,
  store_id uuid NOT NULL REFERENCES public.stores(id) ON DELETE CASCADE,
  actor_user_id uuid NOT NULL REFERENCES public.app_users(id),
  entity_type text NOT NULL,
  entity_id uuid NOT NULL,
  action text NOT NULL,
  before_data jsonb NULL,
  after_data jsonb NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
