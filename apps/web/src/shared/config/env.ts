const DEFAULT_API_BASE_URL = "http://localhost:8080";
const KNOWN_APP_ENVS = ["development", "test", "production"] as const;

type AppEnv = (typeof KNOWN_APP_ENVS)[number];

function isAppEnv(value: string): value is AppEnv {
  return KNOWN_APP_ENVS.includes(value as AppEnv);
}

function readEnvValue(value: unknown): string | undefined {
  if (typeof value !== "string") {
    return undefined;
  }

  return value;
}

function normalizeAppEnv(value: string | undefined): AppEnv {
  if (value && isAppEnv(value)) {
    return value;
  }

  return "development";
}

function normalizeApiBaseUrl(value: string | undefined): string {
  const trimmed = value?.trim();
  const baseUrl = trimmed && trimmed.length > 0 ? trimmed : DEFAULT_API_BASE_URL;
  return baseUrl.replace(/\/$/, "");
}

export const env = {
  appEnv: normalizeAppEnv(readEnvValue(import.meta.env.VITE_APP_ENV)),
  apiBaseUrl: normalizeApiBaseUrl(readEnvValue(import.meta.env.VITE_API_BASE_URL)),
};
