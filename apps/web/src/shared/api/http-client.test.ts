import { afterEach, describe, expect, it, vi } from "vitest";

import { apiRequest } from "./http-client";
import { clearAccessToken, setAccessToken } from "../auth/token-storage";

describe("apiRequest", () => {
  afterEach(() => {
    clearAccessToken();
    vi.restoreAllMocks();
  });

  it("attaches authorization header when token exists", async () => {
    setAccessToken("token-123");

    const fetchMock = vi.fn(async (_input: RequestInfo | URL, init?: RequestInit) => {
      const headers = new Headers(init?.headers);
      expect(headers.get("Authorization")).toBe("Bearer token-123");

      return new Response(JSON.stringify({ status: "ok" }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    });

    vi.stubGlobal("fetch", fetchMock);

    await expect(apiRequest<{ status: string }>("/api/v1/auth/ping")).resolves.toEqual({
      status: "ok",
    });
  });

  it("does not attach authorization when requireAuth is false", async () => {
    setAccessToken("token-123");

    const fetchMock = vi.fn(async (_input: RequestInfo | URL, init?: RequestInit) => {
      const headers = new Headers(init?.headers);
      expect(headers.get("Authorization")).toBeNull();

      return new Response(JSON.stringify({ status: "ok" }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    });

    vi.stubGlobal("fetch", fetchMock);

    await expect(
      apiRequest<{ status: string }>("/api/v1/public", { requireAuth: false }),
    ).resolves.toEqual({ status: "ok" });
  });
});
