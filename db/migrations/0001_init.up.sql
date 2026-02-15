BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE FUNCTION public.set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$;

CREATE TABLE public.app_users (
  id uuid PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
  display_name text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE public.stores (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  created_by uuid NOT NULL REFERENCES public.app_users(id),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE public.store_memberships (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  store_id uuid NOT NULL REFERENCES public.stores(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES public.app_users(id) ON DELETE CASCADE,
  role text NOT NULL CHECK (role IN ('manager', 'staff')),
  status text NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'disabled')),
  joined_at timestamptz NOT NULL DEFAULT now(),
  disabled_at timestamptz NULL,
  UNIQUE (store_id, user_id)
);

CREATE INDEX idx_memberships_user ON public.store_memberships (user_id);
CREATE INDEX idx_memberships_store_role ON public.store_memberships (store_id, role);

CREATE TABLE public.invites (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
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
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
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

CREATE INDEX idx_periods_store_dates ON public.schedule_periods (store_id, start_date DESC);
CREATE INDEX idx_periods_store_status ON public.schedule_periods (store_id, status);

CREATE TABLE public.availability_submissions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  period_id uuid NOT NULL REFERENCES public.schedule_periods(id) ON DELETE CASCADE,
  staff_user_id uuid NOT NULL REFERENCES public.app_users(id),
  submitted_by uuid NOT NULL REFERENCES public.app_users(id),
  submitted_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (period_id, staff_user_id)
);

CREATE INDEX idx_avail_submission_period ON public.availability_submissions (period_id);

CREATE TABLE public.availability_entries (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  submission_id uuid NOT NULL REFERENCES public.availability_submissions(id) ON DELETE CASCADE,
  work_date date NOT NULL,
  start_time time NOT NULL,
  end_time time NOT NULL,
  availability_type text NOT NULL DEFAULT 'available' CHECK (availability_type IN ('available', 'unavailable')),
  CHECK (start_time < end_time)
);

CREATE INDEX idx_avail_entries_submission_date
  ON public.availability_entries (submission_id, work_date);

CREATE TABLE public.shift_assignments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
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

CREATE INDEX idx_assignments_period_date ON public.shift_assignments (period_id, work_date);
CREATE INDEX idx_assignments_staff_date ON public.shift_assignments (staff_user_id, work_date);
CREATE INDEX idx_assignments_store_period ON public.shift_assignments (store_id, period_id);

CREATE TABLE public.coverage_targets (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  period_id uuid NOT NULL REFERENCES public.schedule_periods(id) ON DELETE CASCADE,
  day_of_week int NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
  start_time time NOT NULL,
  end_time time NOT NULL,
  required_count int NOT NULL CHECK (required_count >= 0),
  CHECK (start_time < end_time),
  UNIQUE (period_id, day_of_week, start_time, end_time)
);

CREATE TABLE public.change_requests (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
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

CREATE INDEX idx_change_requests_store_status ON public.change_requests (store_id, status);
CREATE INDEX idx_change_requests_requester
  ON public.change_requests (requester_user_id, created_at DESC);

CREATE TABLE public.notifications (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
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

CREATE INDEX idx_notifications_user_created
  ON public.notifications (user_id, created_at DESC);
CREATE INDEX idx_notifications_unread
  ON public.notifications (user_id)
  WHERE read_at IS NULL;

CREATE TABLE public.audit_logs (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  store_id uuid NOT NULL REFERENCES public.stores(id) ON DELETE CASCADE,
  actor_user_id uuid NOT NULL REFERENCES public.app_users(id),
  entity_type text NOT NULL,
  entity_id uuid NOT NULL,
  action text NOT NULL,
  before_data jsonb NULL,
  after_data jsonb NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_store_created ON public.audit_logs (store_id, created_at DESC);
CREATE INDEX idx_audit_logs_entity ON public.audit_logs (entity_type, entity_id);

CREATE TRIGGER trg_app_users_set_updated_at
BEFORE UPDATE ON public.app_users
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

CREATE TRIGGER trg_stores_set_updated_at
BEFORE UPDATE ON public.stores
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

CREATE TRIGGER trg_schedule_periods_set_updated_at
BEFORE UPDATE ON public.schedule_periods
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

CREATE TRIGGER trg_availability_submissions_set_updated_at
BEFORE UPDATE ON public.availability_submissions
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

CREATE TRIGGER trg_shift_assignments_set_updated_at
BEFORE UPDATE ON public.shift_assignments
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

CREATE TRIGGER trg_change_requests_set_updated_at
BEFORE UPDATE ON public.change_requests
FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

COMMIT;
