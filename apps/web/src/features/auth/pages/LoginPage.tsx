import { type SyntheticEvent, useMemo, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

import { useAuth } from "../hooks/useAuth";

interface LocationState {
  from?: string;
}

function resolveDestination(state: unknown): string {
  if (!state || typeof state !== "object") {
    return "/app";
  }

  const from = (state as LocationState).from;
  if (!from || typeof from !== "string") {
    return "/app";
  }

  return from;
}

export function LoginPage() {
  const [token, setToken] = useState("");
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const destination = useMemo(() => resolveDestination(location.state), [location.state]);

  function handleSubmit(event: SyntheticEvent<HTMLFormElement>) {
    event.preventDefault();

    const trimmedToken = token.trim();
    if (!trimmedToken) {
      setErrorMessage("トークンを入力してください");
      return;
    }

    login(trimmedToken);
    void navigate(destination, { replace: true });
  }

  return (
    <main className="mx-auto flex min-h-screen w-full max-w-md items-center px-6 py-16">
      <section className="w-full rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
        <h1 className="text-2xl font-semibold tracking-tight text-slate-900">ログイン</h1>
        <p className="mt-2 text-sm text-slate-600">
          BL-007時点ではテスト用トークンでログインします。後続でSupabase認証に置き換えます。
        </p>

        <form className="mt-6 space-y-4" onSubmit={handleSubmit}>
          <div className="space-y-2">
            <label className="text-sm font-medium text-slate-800" htmlFor="access-token">
              アクセストークン
            </label>
            <input
              id="access-token"
              autoComplete="off"
              className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm text-slate-900 outline-none transition focus:border-slate-500 focus:ring-2 focus:ring-slate-200"
              placeholder="eyJhbGciOi..."
              type="text"
              value={token}
              onChange={(event) => {
                setToken(event.target.value);
                setErrorMessage(null);
              }}
            />
          </div>

          {errorMessage ? <p className="text-sm text-red-700">{errorMessage}</p> : null}

          <button
            className="inline-flex w-full items-center justify-center rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white transition hover:bg-slate-700"
            type="submit"
          >
            ログイン
          </button>
        </form>
      </section>
    </main>
  );
}
