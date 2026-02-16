import { useAuth, useAuthPingQuery } from "../../auth";

export function DashboardPage() {
  const { logout } = useAuth();
  const pingQuery = useAuthPingQuery();

  return (
    <main className="mx-auto min-h-screen w-full max-w-5xl px-6 py-12">
      <header className="flex flex-col gap-3 border-b border-slate-200 pb-6 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight text-slate-900">Shift Manager</h1>
          <p className="mt-1 text-sm text-slate-600">BL-007 Web骨組み（認証ガード + API接続）</p>
        </div>
        <button
          className="inline-flex items-center justify-center rounded-md border border-slate-300 px-4 py-2 text-sm font-medium text-slate-700 transition hover:bg-slate-100"
          type="button"
          onClick={logout}
        >
          ログアウト
        </button>
      </header>

      <section className="mt-8 grid gap-4 sm:grid-cols-2">
        <article className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h2 className="text-sm font-semibold uppercase tracking-wider text-slate-500">Auth state</h2>
          <p className="mt-3 text-sm text-slate-800">ログイン状態: 認証済み</p>
        </article>

        <article className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <h2 className="text-sm font-semibold uppercase tracking-wider text-slate-500">API ping</h2>
          {pingQuery.isPending ? <p className="mt-3 text-sm text-slate-700">接続確認中...</p> : null}

          {pingQuery.isError ? (
            <p className="mt-3 text-sm text-red-700">API接続に失敗しました。トークンとAPI起動状態を確認してください。</p>
          ) : null}

          {pingQuery.isSuccess ? (
            <p className="mt-3 text-sm text-emerald-700">
              接続OK: {pingQuery.data.module} / {pingQuery.data.status}
            </p>
          ) : null}
        </article>
      </section>
    </main>
  );
}
