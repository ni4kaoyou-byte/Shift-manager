import { getAccessToken } from "../auth/token-storage";
import { env } from "../config/env";

interface ApiErrorPayload {
  error?: {
    code?: string;
    message?: string;
  };
}

type JsonValue = Record<string, unknown> | unknown[];

export interface ApiRequestInit extends Omit<RequestInit, "headers" | "body"> {
  headers?: HeadersInit;
  body?: JsonValue;
  requireAuth?: boolean;
}

export class ApiClientError extends Error {
  public readonly status: number;
  public readonly code: string;

  public constructor(status: number, code: string, message: string) {
    super(message);
    this.name = "ApiClientError";
    this.status = status;
    this.code = code;
  }
}

function toHeaders(headers?: HeadersInit): Headers {
  return new Headers(headers);
}

async function parseJson(response: Response): Promise<unknown> {
  const text = await response.text();
  if (!text) {
    return null;
  }

  return JSON.parse(text) as unknown;
}

function getErrorPayload(value: unknown): ApiErrorPayload {
  if (value && typeof value === "object") {
    return value as ApiErrorPayload;
  }

  return {};
}

export async function apiRequest<TResponse>(
  path: string,
  init: ApiRequestInit = {},
): Promise<TResponse> {
  const { requireAuth = true, body, headers, ...requestInit } = init;

  const requestHeaders = toHeaders(headers);
  requestHeaders.set("Accept", "application/json");

  if (body) {
    requestHeaders.set("Content-Type", "application/json");
  }

  if (requireAuth) {
    const token = getAccessToken();
    if (token) {
      requestHeaders.set("Authorization", `Bearer ${token}`);
    }
  }

  const response = await fetch(`${env.apiBaseUrl}${path}`, {
    ...requestInit,
    headers: requestHeaders,
    body: body ? JSON.stringify(body) : undefined,
  });

  const payload = await parseJson(response);

  if (!response.ok) {
    const errorPayload = getErrorPayload(payload);
    const code = errorPayload.error?.code ?? `http_${response.status}`;
    const message = errorPayload.error?.message ?? "request failed";
    throw new ApiClientError(response.status, code, message);
  }

  return payload as TResponse;
}
