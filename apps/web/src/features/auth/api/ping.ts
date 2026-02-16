import { apiRequest } from "../../../shared/api/http-client";

export interface AuthPingResponse {
  module: string;
  status: string;
}

export function fetchAuthPing(): Promise<AuthPingResponse> {
  return apiRequest<AuthPingResponse>("/api/v1/auth/ping", {
    method: "GET",
  });
}
