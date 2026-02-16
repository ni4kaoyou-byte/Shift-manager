import { type PropsWithChildren, useCallback, useMemo, useState } from "react";

import { clearAccessToken, getAccessToken, setAccessToken } from "../../../shared/auth/token-storage";
import { AuthContext } from "./auth-context";

export function AuthProvider({ children }: PropsWithChildren) {
  const [accessToken, setAccessTokenState] = useState<string | null>(() => getAccessToken());

  const login = useCallback((token: string) => {
    const trimmedToken = token.trim();
    setAccessToken(trimmedToken);
    setAccessTokenState(trimmedToken);
  }, []);

  const logout = useCallback(() => {
    clearAccessToken();
    setAccessTokenState(null);
  }, []);

  const value = useMemo(
    () => ({
      accessToken,
      isAuthenticated: Boolean(accessToken),
      login,
      logout,
    }),
    [accessToken, login, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
