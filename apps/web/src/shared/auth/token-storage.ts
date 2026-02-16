let accessToken: string | null = null;

export function getAccessToken(): string | null {
  return accessToken;
}

export function setAccessToken(token: string): void {
  const normalized = token.trim();
  accessToken = normalized.length > 0 ? normalized : null;
}

export function clearAccessToken(): void {
  accessToken = null;
}
