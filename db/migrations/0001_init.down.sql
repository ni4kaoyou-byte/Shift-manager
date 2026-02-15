BEGIN;

DROP TABLE IF EXISTS public.audit_logs;
DROP TABLE IF EXISTS public.notifications;
DROP TABLE IF EXISTS public.change_requests;
DROP TABLE IF EXISTS public.coverage_targets;
DROP TABLE IF EXISTS public.shift_assignments;
DROP TABLE IF EXISTS public.availability_entries;
DROP TABLE IF EXISTS public.availability_submissions;
DROP TABLE IF EXISTS public.schedule_periods;
DROP TABLE IF EXISTS public.invites;
DROP TABLE IF EXISTS public.store_memberships;
DROP TABLE IF EXISTS public.stores;
DROP TABLE IF EXISTS public.app_users;

DROP FUNCTION IF EXISTS public.set_updated_at();

COMMIT;
