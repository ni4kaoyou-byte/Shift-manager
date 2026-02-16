import { Navigate, Outlet } from "react-router-dom";

import { useAuth } from "../hooks/useAuth";

export function PublicOnlyRoute() {
  const { isAuthenticated } = useAuth();

  if (isAuthenticated) {
    return <Navigate to="/app" replace />;
  }

  return <Outlet />;
}
