import { Link, Route, Routes } from "react-router-dom";

function HomePage() {
  return (
    <main className="mx-auto max-w-3xl px-6 py-16">
      <h1 className="text-3xl font-bold">Shift Manager</h1>
      <p className="mt-3 text-gray-600">
        React + TypeScript + Vite の初期セットアップが完了しています。
      </p>
      <div className="mt-8">
        <Link className="text-blue-600 underline" to="/health">
          ヘルス画面へ
        </Link>
      </div>
    </main>
  );
}

function HealthPage() {
  return (
    <main className="mx-auto max-w-3xl px-6 py-16">
      <h2 className="text-2xl font-semibold">Web Health</h2>
      <p className="mt-3 text-gray-600">Frontend scaffolding is ready.</p>
      <div className="mt-6">
        <Link className="text-blue-600 underline" to="/">
          トップに戻る
        </Link>
      </div>
    </main>
  );
}

export function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/health" element={<HealthPage />} />
    </Routes>
  );
}
