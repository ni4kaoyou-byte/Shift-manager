import { Link } from "react-router-dom";

export function NotFoundPage() {
  return (
    <main className="mx-auto min-h-screen w-full max-w-3xl px-6 py-16">
      <h1 className="text-2xl font-semibold text-slate-900">ページが見つかりません</h1>
      <p className="mt-2 text-sm text-slate-600">URLを確認するか、ダッシュボードに戻ってください。</p>
      <Link className="mt-6 inline-flex text-sm font-medium text-slate-900 underline" to="/app">
        ダッシュボードへ戻る
      </Link>
    </main>
  );
}
