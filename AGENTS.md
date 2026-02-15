# AGENTS.md

このリポジトリで作業するCodex向け運用ルール。

## 目的
- Go / React の lint・format・test 品質を常に維持する。
- 変更後に「ローカルでは通るがCIで落ちる」を防ぐ。

## 必須チェック（コード変更時）

### API（`apps/api`）
1. `make fmt`
2. `make lint`
3. `make test`

### Web（`apps/web`）
1. `npm run format`
2. `npm run lint`
3. `npm run test`
4. `npm run build`

## 実行タイミング
- コード変更を含むタスクの完了前に必ず実行する。
- 失敗した場合は修正して再実行する。
- 失敗が解消できない場合は、失敗コマンドと原因を明示して報告する。

## 運用メモ
- Go lint は `golangci-lint` を正とする（設定: `apps/api/.golangci.yml`）。
- Web lint は `eslint.config.js` を正とする。
- Web format は `prettier`（`apps/web/.prettierrc.json`）を正とする。
