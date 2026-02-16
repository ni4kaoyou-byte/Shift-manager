# 11_react-review-checklist.md — React feature-based レビュー・チェックリスト

対象: `apps/web`

---

## 1) ディレクトリ設計 / 責務の境界

- [ ] 変更は関係する feature 配下に閉じている（別 feature への「ついで修正」を避ける）
- [ ] `pages/` は画面の合成（composition）中心で、ロジックや API 呼び出しを抱えすぎていない
- [ ] `components/` は UI 表示に集中している（副作用・API・グローバル状態に依存しない）
- [ ] `hooks/` は状態と副作用のカプセル化になっている（UI 詳細を持ち込まない）
- [ ] `api/` は通信だけ（UI整形/状態管理の責務を持たない）
- [ ] feature 間の import が片方向またはルールに沿っている（相互依存が増えていない）

あるあるNG:
- `src/features/auth/hooks/useAuth.ts` が schedule の API を直接叩く
- 「auth の関心事」ではない処理が混ざるのは境界崩れのサイン

---

## 2) 依存方向 / import ルール

- [ ] feature 外への依存は共通層（`shared/` や `lib/`）に寄せられている
- [ ] 別 feature の `components/` / `hooks/` を直接 import していない（必要なら公開APIを作る）
- [ ] `pages -> components/hooks/api` の方向はOK、逆向きがない
- [ ] import path が一貫（例: `@/features/...`）し、相対パス地獄になっていない

おすすめ運用:
- feature を跨ぐ共有は `shared/ui`, `shared/lib`, `shared/api` に置く
- 他 feature の内部を直接触らないルールにする

---

## 3) API 通信（api/）の品質

- [ ] API関数は入出力が型で固定されている（`any` を混ぜない）
- [ ] エラーを握りつぶしていない（呼び出し側が判断できる形で返す/throwする）
- [ ] リクエスト責務が明確（UI都合の整形を `api/` が抱えすぎない）
- [ ] キャッシュ戦略が一貫（SWR/React Query のキー設計・invalidate が妥当）
- [ ] `api/` が `toast` や `router` を直接触っていない（副作用は上位へ）

---

## 4) 状態管理 / hooks のレビュー観点

- [ ] 状態の置き場所が妥当（ローカル / feature内 / グローバル）
- [ ] `useEffect` が何でも屋になっていない（依存配列が正しい、無限ループしない）
- [ ] 同期・非同期状態が分離されている（`loading/error/data` の扱いが破綻していない）
- [ ] 派生 state を `useState` に持っていない（`useMemo` 等で計算できるものは分離）
- [ ] hooks の戻り値が使いやすく統一されている（例: `{ data, isLoading, error, actions }`）

---

## 5) UI / components の観点

- [ ] コンポーネントが大きすぎない（責務分割されていて読める）
- [ ] props の過剰なバケツリレーがない（必要なら composition / context を検討）
- [ ] 表示とデータ取得が混ざっていない（Container/Presenter 的な分離）
- [ ] ローディング/空/エラー状態が実装されている（無言で固まらない）
- [ ] a11y 最低限（`button`/`div` の使い分け、aria、フォーカス、キーボード操作）

---

## 6) 仕様の正しさ（feature 単位）

- [ ] 業務ルールが UI に散らばっていない（hooks / domain的関数に寄せる）
- [ ] 権限/認可/表示制御が一貫している（見えるだけ・押せるだけの抜けがない）
- [ ] 日付・時刻・タイムゾーン（シフト系の地雷）を正しく扱っている

---

## 7) テスト / 回帰リスク

- [ ] 重要ロジックがテスト可能な形になっている（UIに埋め込まない）
- [ ] 壊れやすい箇所に最低限のテストがある（バリデーション、状態遷移、API成功/失敗）
- [ ] モック戦略が破綻していない（API層をモックできる構造）

---

## 8) パフォーマンス / UX

- [ ] 不要な再レンダが目立たない（重い計算、無駄な state、`key` ミスの回避）
- [ ] 大量リスト対策がある（ページング/仮想化/スケルトン）
- [ ] optimistic update / refetch UX が妥当（チラつき、二重送信がない）

---

## 9) セキュリティ / 安全性

- [ ] ユーザー入力が危険な出力（HTML挿入等）に直結していない
- [ ] トークン等の秘匿情報をログ出力していない
- [ ] APIエラーをそのまま露出していない（表示文言が適切）

---

## 10) メンテ性 / PR 品質

- [ ] 命名が一貫（`changeRequest` vs `change-request` などが揺れていない）
- [ ] 「ついでリファクタ」が多すぎず、レビュースコープが破裂していない
- [ ] PR説明で変更理由が追える（何を・なぜ・どう確認したか）
- [ ] TODO/コメントが放置されていない（残すなら期限/理由を書く）

---

## PRテンプレ（貼り付け用）

- 変更したfeature: auth / schedule / change-request / ...
- 影響範囲: 画面 / API / 状態 / ルーティング / 権限
- 確認方法: 手動手順 / スクショ / テスト結果
- 懸念点: 回帰しそうな点、境界が怪しい点

---

## 追加ルール（feature-basedで特に効く）

- 跨ぎ import 禁止:
  - 他 feature の内部（`components/hooks/api`）を直接参照しない
  - 必要なら `features/<x>/index.ts` などの公開口を作る
- `pages` は合成だけ:
  - API直叩きや巨大 `useEffect` を置かない
- `api` は副作用ゼロ:
  - `toast` / `router` / DOM操作は禁止
