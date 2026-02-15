import js from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import eslintConfigPrettier from "eslint-config-prettier";

import unusedImports from "eslint-plugin-unused-imports";
import importX from "eslint-plugin-import-x";
import jsxA11y from "eslint-plugin-jsx-a11y";
import simpleImportSort from "eslint-plugin-simple-import-sort";

const isTestFile = ["**/*.{test,spec}.{ts,tsx}", "**/__tests__/**", "**/tests/**"];

export default tseslint.config(
  // 1) ignore
  { ignores: ["dist", "coverage", "node_modules"] },

  // 2) App (strict)
  {
    files: ["**/*.{ts,tsx}"],
    ignores: isTestFile,

    extends: [
      js.configs.recommended,

      // 🔥 Type-aware (最重要)
      ...tseslint.configs.recommendedTypeChecked,
      ...tseslint.configs.strictTypeChecked,
      ...tseslint.configs.stylisticTypeChecked,

      eslintConfigPrettier,
    ],

    languageOptions: {
      ecmaVersion: 2022,
      globals: globals.browser,
      parserOptions: {
        // 型情報ありルールを動かすやつ（TS5+ならこれが楽）
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },

    plugins: {
      "react-hooks": reactHooks,
      "react-refresh": reactRefresh,

      "unused-imports": unusedImports,
      "import-x": importX,
      "jsx-a11y": jsxA11y,
      "simple-import-sort": simpleImportSort,
    },

    linterOptions: {
      reportUnusedDisableDirectives: true,
    },

    rules: {
      /* -------------------------
       * React
       * ------------------------- */
      ...reactHooks.configs.recommended.rules,
      "react-hooks/exhaustive-deps": "error",
      "react-refresh/only-export-components": ["error", { allowConstantExport: true }],

      /* -------------------------
       * Imports / Sorting
       * ------------------------- */
      "unused-imports/no-unused-imports": "error",
      "unused-imports/no-unused-vars": [
        "error",
        {
          vars: "all",
          varsIgnorePattern: "^_",
          args: "after-used",
          argsIgnorePattern: "^_",
        },
      ],
      "simple-import-sort/imports": "error",
      "simple-import-sort/exports": "error",

      // importの健全性（解決できない/重複など）
      "import-x/no-duplicates": "error",
      "import-x/no-mutable-exports": "error",
      "import-x/first": "error",
      "import-x/newline-after-import": "error",

      /* -------------------------
       * a11y（UIの地雷を先に潰す）
       * ------------------------- */
      ...jsxA11y.configs.recommended.rules,

      /* -------------------------
       * TS “事故防止” 本丸
       * ------------------------- */
      "@typescript-eslint/consistent-type-imports": ["error", { prefer: "type-imports" }],
      "@typescript-eslint/no-floating-promises": "error",
      "@typescript-eslint/no-misused-promises": [
        "error",
        { checksVoidReturn: { attributes: false } },
      ],
      "@typescript-eslint/no-unnecessary-condition": "error",
      "@typescript-eslint/no-unnecessary-type-assertion": "error",
      "@typescript-eslint/no-unsafe-argument": "error",
      "@typescript-eslint/no-unsafe-assignment": "error",
      "@typescript-eslint/no-unsafe-call": "error",
      "@typescript-eslint/no-unsafe-member-access": "error",
      "@typescript-eslint/no-unsafe-return": "error",
      "@typescript-eslint/restrict-template-expressions": [
        "error",
        { allowNumber: true, allowBoolean: true, allowNullish: true },
      ],
      "@typescript-eslint/switch-exhaustiveness-check": "error",

      // anyは「全面禁止」にしてもいい（きついけど最強）
      "@typescript-eslint/no-explicit-any": "error",

      /* -------------------------
       * JS/General
       * ------------------------- */
      eqeqeq: ["error", "always"],
      "no-debugger": "error",
      "no-console": ["error", { allow: ["warn", "error"] }],
    },
  },

  // 3) Tests (loose)
  {
    files: isTestFile,
    extends: [
      js.configs.recommended,
      ...tseslint.configs.recommended, // ← type-aware外して軽く
      eslintConfigPrettier,
    ],
    languageOptions: {
      ecmaVersion: 2022,
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
    plugins: {
      "unused-imports": unusedImports,
    },
    rules: {
      // テストはガンガン書けるように緩める
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/no-non-null-assertion": "off",
      "@typescript-eslint/no-unsafe-assignment": "off",
      "@typescript-eslint/no-unsafe-call": "off",
      "@typescript-eslint/no-unsafe-member-access": "off",
      "@typescript-eslint/no-unsafe-return": "off",
      "@typescript-eslint/no-floating-promises": "off",
      "@typescript-eslint/no-misused-promises": "off",
      "@typescript-eslint/no-unnecessary-condition": "off",

      "no-console": "off",
      "unused-imports/no-unused-imports": "warn",
    },
  },
);
